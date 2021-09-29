"use strict";

function mobileNavToggle() {
    var menu = document.getElementById('mobile-menu').parentElement;
    menu.classList.toggle('mobile-menu-visible');
}

function docsVersionToggle() {
    var menu = document.getElementById('dropdown-menu');
    menu.classList.toggle('dropdown-menu-visible');
}

window.onclick = function(event) {
    var 
        target = event.target,
        menu = document.getElementById('dropdown-menu')
    ;

    if(menu !== null) {
        if(!target.classList.contains('dropdown-toggle')) {
            menu.classList.remove('dropdown-menu-visible');
        }
    }

}

function toggleAccordion(el) {
    var state = el.getAttribute('aria-expanded');
            
    // toggle state
    state = (state === 'true') ? 'false' : 'true';
    el.setAttribute('aria-expanded', state);

    // show/hide list
    el.nextElementSibling.classList.toggle('show');
}

// Adds copy to clipboard buttons to all codeblocks for hugo rendered site
// Ref: https://github.com/dguo/dannyguo.com/blob/main/content/blog/how-to-add-copy-to-clipboard-buttons-to-code-blocks-in-hugo.md
function addCopyButtons(clipboard) {
    document.querySelectorAll('pre > code').forEach(function (codeBlock) {
        var button = document.createElement('button');
        button.className = 'copy-code-button';
        button.type = 'button';
        button.innerText = 'Copy';

        button.addEventListener('click', function () {
            clipboard.writeText(codeBlock.innerText.trim()).then(function () {
                /* Chrome doesn't seem to blur automatically,
                   leaving the button in a focused state. */
                button.blur();
                button.innerText = 'Copied!';
                button.classList.add('copy-code-button-copied')

                setTimeout(function () {
                    button.innerText = 'Copy';
                    button.classList.remove('copy-code-button-copied');
                }, 500);
            }, function (error) {
                button.innerText = 'Error!';
                button.classList.add('copy-code-button-error')
                console.error("could not copy to clipboard");
                console.error(error)

                setTimeout(function () {
                    button.innerText = 'Copy';
                    button.classList.remove('copy-code-button-error');
                }, 500);
            });
        });

        var pre = codeBlock.parentNode;
        if (pre.parentNode.classList.contains('highlight')) {
            var highlight = pre.parentNode;
            highlight.parentNode.insertBefore(button, highlight);
        } else {
            pre.parentNode.insertBefore(button, pre);
        }
    });
}

function createCopyButtons() {
    if (navigator && navigator.clipboard) {
        addCopyButtons(navigator.clipboard);
    } else {
        // navigator.clipboard is supported in all modern browsers except for Internet Explorer
        // https://developer.mozilla.org/en-US/docs/Web/API/Navigator/clipboard#browser_compatibility
        console.warn("Code copy buttons not supported in browser!");
    }
}

document.addEventListener('DOMContentLoaded', function(){
    // hamburger
    var hamburger = document.getElementById('mobileNavToggle');
    hamburger.addEventListener('click', function() {
        mobileNavToggle();
    });

    // accordion
    var collapsible = document.getElementsByClassName('collapse-trigger');
    
    if(collapsible.length) {
        for (var i = 0; i < collapsible.length; i++) {
            collapsible[i].addEventListener('click', function() {
                toggleAccordion(this);
            });
        }
        
        // open accordion for active doc section
        var active = document.querySelectorAll('.collapse .active');
        
        if(active.length) {
            toggleAccordion(active[0].closest('.collapse').previousElementSibling);
        } else {
            toggleAccordion(document.querySelector('.collapse-trigger'));
            document.querySelector('.collapse a').classList.add('active');
        }
    }

    var dropdown = document.getElementById('dropdownMenuButton');
    
    if(dropdown !== null) {
        dropdown.addEventListener('click', function() {
            docsVersionToggle();
        });
    }

    createCopyButtons();

    // Load the medium-zoom library and attach based on css selector
    mediumZoom('.docs img')
});
