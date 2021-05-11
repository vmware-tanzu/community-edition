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

    if(!target.classList.contains('dropdown-toggle')) {
        menu.classList.remove('dropdown-menu-visible');
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

document.addEventListener('DOMContentLoaded', function(){
    // accordion
    var collapsible = document.getElementsByClassName('collapse-trigger');

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
});
