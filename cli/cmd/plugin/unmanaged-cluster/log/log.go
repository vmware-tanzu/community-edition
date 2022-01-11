// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package log provides logging mechanisms for the tanzu unmanaged-cluster CLI plugin. It offers
// logging functionality that can include stylized logs, updating progress dots (...), and emojis.
// It also respects a TTY parameter. When set to false, all stylization is removed.
package log

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

const (
	// The following are a set of emoji codes that can be used
	// with the Event and Eventf logging methods in this package.
	// They should all take up 2 terminal columns.
	// https://www.unicode.org/emoji/charts/full-emoji-list.html
	WrenchEmoji     = "\\U+1F527"
	FolderEmoji     = "\\U+1F4C1"
	PictureEmoji    = "\\U+1F3A8"
	PackageEmoji    = "\\U+1F4E6"
	RocketEmoji     = "\\U+1F680"
	EnvelopeEmoji   = "\\U+1F4E7"
	GlobeEmoji      = "\\U+1F310"
	GreenCheckEmoji = "\\U+2705"
	ControllerEmoji = "\\U+1F3AE"
	TestTubeEmoji   = "\\U+1F9EA"
)

// CMDLogger is the logger implementation used for high-level command line logging.
type CMDLogger struct {
	// whether to support stylizing logging output
	tty bool
	// logging level to respect for this logger
	level int
	// log level set by a logging event
	logLevel int
	// instances of indentation (" ") to prepend to a long message
	indent int
	// logColor defines the color to log the message as define by fatih/color Attributes
	logColor color.Attribute
	// output controls where log messages are sent
	output io.Writer
}

// Logger provides the logging interaction for the application.
type Logger interface {
	// Event takes an emoji codepoint (e.g. "\\U+1F609") and presents a log message.
	// This package provides several standard emoji codepoints as string constants. I.e: logger.HammerEmoji
	// Warning: Emojis may have variable width and this method assumes 2 width emojis, adding a space between the emoji and message.
	// Emojis provided in this package as string consts have 2 width and work with this method.
	// If you wish for additional space, add it at the beginning of the message (string) argument.
	Event(emoji, message string)
	// Eventf takes an emoji codepoint (e.g. "\\U+1F609"), a format string, arguments for the format string.
	// This package provides several standard emoji codepoints as string constants. I.e: logger.HammerEmoji
	// Warning: Emojis may have variable width and this method assumes 2 width emojis, adding a space between the emoji and message.
	// Emojis provided in this package as string consts have 2 width and work with this method.
	// If you wish for additional space, add it at the beginning of the message (string) argument.
	Eventf(emoji, message string, args ...interface{})
	// Info prints a standard log message.
	// Line breaks are not automatically added to the end.
	Info(message string)
	// Infof takes a format string, arguments, and prints it as a standard log message.
	// Line breaks are not automatically added to the end.
	Infof(message string, args ...interface{})
	// Warn prints a warning message. When TTY is enabled (default), it will be stylized as yellow.
	// Line breaks are not automatically added to the end.
	Warn(message string)
	// Warnf takes a format string, arguments, and prints it as a warning message.
	// When TTY is enabled (default), it will be stylized as yellow.
	// Line breaks are not automatically added to the end.
	Warnf(message string, args ...interface{})
	// Error prints an error message. When TTY is enabled (default), it will be stylized as red.
	// Line breaks are not automatically added to the end.
	Error(message string)
	// Errorf takes a format string, arguments, and prints it as an error message.
	// When TTY is enabled (default), it will be stylized as yellow.
	// Line breaks are not automatically added to the end.
	Errorf(message string, args ...interface{})
	// ReplaceLinef takes a template string message
	// and any optional format arguments
	// and replaces the current line.
	// This is useful after canceling AnimateProgressWithOptions and needing to print a final "success" message
	// Ex: ReplaceLinef("Finished reconciling controller: %s", controllerStatus)
	ReplaceLinef(message string, args ...interface{})
	// AnimateProgressWithOptions takes any number of AnimatorOptions
	// and is used to async animate a number of dots.
	// See the AnimatorOptions for further documentation
	// Ex: AnimateProgressWithOptions(AnimatorWithMaxLen(5))
	AnimateProgressWithOptions(options ...AnimatorOption)
	// V sets the level of the log message based on an integer. The logger implementation will hold a configured
	// log level, which this V level is assessed against to determine whether the log message should be output.
	V(level int) Logger
	// Style provides indentation and colorization of log messages. The indent argument specifies the amount of " "
	// characters to prepend to the message. The color should be specified using color constants in this package.
	Style(indent int, c color.Attribute) Logger
	// AddLogFile adds a file name to log all activity to.
	AddLogFile(filePath string)
}

// NewLogger returns an instance of Logger, implemented via CMDLogger.
func NewLogger(tty bool, level int) Logger {
	return &CMDLogger{
		tty:    tty,
		level:  level,
		output: os.Stdout,
	}
}

func (l *CMDLogger) AddLogFile(filePath string) {
	logFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		l.Warnf("Failed to open log file %q: %v", filePath, err)
		return
	}

	l.output = io.MultiWriter(logFile, os.Stdout)
}

func (l *CMDLogger) Event(emoji, message string) {
	if l.logLevel > l.level {
		return
	}
	// when tty is off, remove emoji from output
	if !l.tty {
		emoji = ""
		// space is sometimes added to the beginning so that text isn't up against the emoji
		// this trims leading space in case that was one.
		message = strings.TrimLeft(message, " ")
	} else {
		emoji = unquoteCodePoint(emoji)
	}

	// Print a new line before the event is logged
	// so that each event is within it's own "block"
	fmt.Print("\n")

	// process indentation and ensure a space after the emoji
	message = "%s " + message
	message = processStyle(l, message)
	fmt.Fprintf(l.output, message, emoji)
}

func (l *CMDLogger) Eventf(emoji, message string, args ...interface{}) {
	if l.logLevel > l.level {
		return
	}
	// when tty is off, remove emoji from output
	if !l.tty {
		emoji = ""
		// space is sometimes added to the beginning so that text isn't up against the emoji
		// this trims leading space in case that was one.
		message = strings.TrimLeft(message, " ")
	} else {
		emoji = unquoteCodePoint(emoji)
	}

	// Print a new line before the event is logged
	// so that each event is within it's own "block"
	fmt.Print("\n")

	// ensure a space between the emoji and the message
	message = emoji + " " + message
	message = processStyle(l, message)
	fmt.Fprintf(l.output, message, args...)
}

func (l *CMDLogger) Warn(message string) {
	if l.logLevel > l.level {
		return
	}

	message = processStyle(l, message)
	fmt.Print(message)
}

func (l *CMDLogger) Warnf(message string, args ...interface{}) {
	if l.logLevel > l.level {
		return
	}

	message = processStyle(l, message)
	fmt.Fprintf(l.output, message, args...)
}

func (l *CMDLogger) Error(message string) {
	if l.logLevel > l.level {
		return
	}

	message = processStyle(l, message)
	fmt.Print(message)
}

func (l *CMDLogger) Errorf(message string, args ...interface{}) {
	if l.logLevel > l.level {
		return
	}

	message = processStyle(l, message)
	fmt.Fprintf(l.output, message, args...)
}

func (l *CMDLogger) Info(message string) {
	if l.logLevel > l.level {
		return
	}

	message = processStyle(l, message)
	fmt.Print(message)
}

func (l *CMDLogger) Infof(message string, args ...interface{}) {
	if l.logLevel > l.level {
		return
	}

	message = processStyle(l, message)
	fmt.Fprintf(l.output, message, args...)
}

// progressf is an internal method used to log out a specified number of dots
// in addition to a provided message and any format string arguments
func (l *CMDLogger) progressf(count int, message string, args ...interface{}) {
	if l.logLevel > l.level {
		return
	}

	if !l.tty {
		count = 0
	}

	// Add dots to message
	for i := 0; i < count; i++ {
		message += "."
	}

	// Process message style and ensure we clear the line with \r in tty mode
	message = processStyle(l, message)
	if l.tty {
		message = "\r" + message
	}

	// TODO(joshrosso): Is there a better way to do this?
	// we pad with extra space to ensure the line we overwrite (\r) is cleaned
	// nolint
	message += "             "

	// when count is 0, a line break should be added at the end
	// this support non-tty use cases
	if count == 0 {
		message += "\n"
	}

	fmt.Fprintf(l.output, message, args...)
}

func (l *CMDLogger) ReplaceLinef(message string, args ...interface{}) {
	if l.logLevel > l.level {
		return
	}

	// Process message style and Ensure we clear the line with \r in tty mode
	message = processStyle(l, message)
	if l.tty {
		message = "\r" + message
	}

	// TODO(joshrosso): Is there a better way to do this?
	// we pad with extra space to ensure the line we overwrite (\r) is cleaned
	message += "             "

	// add a line break
	// this also supports non-tty use cases
	message += "\n"

	fmt.Fprintf(l.output, message, args...)
}

func (l *CMDLogger) AnimateProgressWithOptions(options ...AnimatorOption) {
	opts := &progressAnimatorOptions{
		maxLen: 5,
	}

	// Apply given animation options
	for _, o := range options {
		o.apply(opts)
	}

	currentLen := 1
	status := ""
	for {
		select {
		case <-opts.ctx.Done():
			return
		case status = <-opts.statChan:
			// noop - this gets the newest status from the status channel
		default:
			// noop - this is used to fallthrough to the processing logic below
			// when there is no status channel or there's no status update
		}

		// Build the format args that eventually get passed to fmt.Fprintf
		// Always expect the status to be first
		fArgs := make([]interface{}, 0)
		if opts.statChan != nil {
			fArgs = append(fArgs, status)
		}

		if len(opts.messagefArgs) != 0 {
			for _, arg := range opts.messagefArgs {
				fArgs = append(fArgs, arg)
			}
		}

		if len(fArgs) == 0 {
			l.progressf(currentLen, opts.messagef)
		} else {
			l.progressf(currentLen, opts.messagef, fArgs...)
		}

		currentLen++
		time.Sleep(1 * time.Second)
		if currentLen == opts.maxLen {
			currentLen = 1
		}
	}
}

func (l *CMDLogger) V(level int) Logger {
	return &CMDLogger{
		tty:      l.tty,
		level:    l.level,
		logLevel: level,
		output:   l.output,
	}
}

func (l *CMDLogger) Style(indent int, c color.Attribute) Logger {
	// if tty is disable, don't return a style-capable logger
	if !l.tty {
		return l
	}
	return &CMDLogger{
		tty:      l.tty,
		level:    l.level,
		logLevel: l.logLevel,
		indent:   indent,
		logColor: c,
		output:   l.output,
	}
}

// unquoteCodePoint takes the unicode value of a symbol and makes it usable for printing.
func unquoteCodePoint(s string) string {
	r, _ := strconv.ParseInt(strings.TrimPrefix(s, "\\U"), 16, 32)
	return string(r)
}

// processStyle adds indentation and color based on the configured CMDLogger. When tty is false, stylization arguments
// are ignored.
func processStyle(l *CMDLogger, message string) string {
	// when tty is off, do no stylization
	if !l.tty {
		return message
	}

	// render indentation
	for i := 0; i < l.indent; i++ {
		message = " " + message
	}

	// apply color value to entire message
	if l.logColor != 0 {
		message = color.New(l.logColor).Sprint(message)
	}

	return message
}
