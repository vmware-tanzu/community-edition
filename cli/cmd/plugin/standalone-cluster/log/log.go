// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package log provides logging mechanisms for the tanzu local CLI plugin. It offers logging functionality that
// can include stylized logs, updating progress dots (...), and emojis. It also respects a TTY parameter. When set to
// false, all stylization is removed.
package log

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	ColorRed        = "\033[31m"
	ColorBlue       = "\033[34m"
	ColorLightGreen = "\033[32m"
	ColorLightGrey  = "\033[37m"
	colorReset      = "\033[0m"
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
	// color to log the message as
	color string
}

// Logger provides the logging interaction for the application.
type Logger interface {
	// Event takes an emoji codepoint (e.g. "\\U+1F609") and presents a log message.
	// Emojis may have variable width, so no space is automatically placed between the emoji and the message.
	// If you wish for additional space, add it at the beginning of the message (string) argument.
	Event(emoji, message string)
	// Eventf takes an emoji codepoint (e.g. "\\U+1F609"), a format string, arguments for the format string.
	// Emojis may have variable width, so no space is automatically placed between the emoji and the message.
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
	// Progressf takes a progress counter, format string, arguments, and prints it as a standard log message.
	// The progress counter will render as a quantity of "." characters defined by the value of count.
	// Usage of Progressf typically involves a for loop feeding in a series of numbers such as
	// 1,2,3,1,2,3,1,2,{exit due to condition}.
	// Progressf will start each message with '\r', which will overwrite the last line. This gives the appearance of
	// updating.
	Progressf(count int, message string, args ...interface{})
	// V sets the level of the log message based on an integer. The logger implementation will hold a configured
	// log level, which this V level is assessed against to determine whether the log message should be output.
	V(level int) Logger
	// Style provides indentation and colorization of log messages. The indent argument specifies the amount of " "
	// characters to prepend to the message. The color should be specified using color constants in this package.
	Style(indent int, color string) Logger
}

// NewLogger returns an instance of Logger, implemented via CMDLogger.
func NewLogger(tty bool, level int) Logger {
	return &CMDLogger{
		tty:   tty,
		level: level,
	}
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

	// process indentation
	message = "%s" + message
	message = processStyle(l, message)
	fmt.Printf(message, emoji)
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
	message = emoji + message
	message = processStyle(l, message)
	fmt.Printf(message, args...)
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
	fmt.Printf(message, args...)
}

func (l *CMDLogger) Error(message string) {
	if l.logLevel > l.level {
		return
	}

	// default to red when color is not set
	if l.color == "" {
		l.color = ColorRed
		message = processStyle(l, message)
	} else {
		message = processStyle(l, message)
	}
	fmt.Print(message)
}

func (l *CMDLogger) Errorf(message string, args ...interface{}) {
	if l.logLevel > l.level {
		return
	}

	// default to red when color is not set
	if l.color == "" {
		l.color = ColorRed
		message = processStyle(l, message)
	} else {
		message = processStyle(l, message)
	}
	message = processStyle(l, message)
	fmt.Printf(message, args...)
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
	fmt.Printf(message, args...)
}

func (l *CMDLogger) Progressf(count int, message string, args ...interface{}) {
	if l.logLevel > l.level {
		return
	}
	if !l.tty {
		count = 0
	}

	for i := 0; i < count; i++ {
		message += "."
	}
	message = processStyle(l, message)
	if l.tty {
		message = "\r" + message
	}
	// TODO(joshrosso): Is there a better way to do this?
	// we pad with extra space to ensure the line we overwrite (\r) is cleaned
	message += "             "
	// when count is 0, a line break should be added at the end
	// this support non-tty use cases
	if count == 0 {
		message += "\n"
	}
	fmt.Printf(message, args...)
}

func (l *CMDLogger) V(level int) Logger {
	return &CMDLogger{
		tty:      l.tty,
		level:    l.level,
		logLevel: level,
	}
}

func (l *CMDLogger) Style(indent int, color string) Logger {
	// if tty is disable, don't return a style-capable logger
	if !l.tty {
		return l
	}
	return &CMDLogger{
		tty:      l.tty,
		level:    l.level,
		logLevel: l.logLevel,
		indent:   indent,
		color:    color,
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
	if l.color != "" {
		message = l.color + message
		message += colorReset
	}

	return message
}
