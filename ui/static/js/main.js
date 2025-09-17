var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
  var link = navLinks[i];
  // Skip the logo link (has nav-logo class) to avoid triangle positioning issues
  if (link.classList.contains("nav-logo")) {
    continue;
  }
  if (link.getAttribute("href") == window.location.pathname) {
    link.classList.add("live");
    break;
  }
}

// Add event listeners for dialog buttons if they exist
var dialogCloseButtons = document.querySelectorAll(".dialog-close-button");
for (var i = 0; i < dialogCloseButtons.length; i++) {
  dialogCloseButtons[i].addEventListener("click", closeDialog);
}

var dialogOpenButtons = document.querySelectorAll(".dialog-open-button");
for (var i = 0; i < dialogOpenButtons.length; i++) {
  dialogOpenButtons[i].addEventListener("click", showDialog);
}

function showDialog() {
  document.querySelector("#dialog-overlay").style.display = "flex";
  document.querySelector("#dialog").style.display = "block";
}

function closeDialog() {
  document.querySelector("#dialog").style.display = "none";
  document.querySelector("#dialog-overlay").style.display = "none";
}

// Hamburger menu functionality
var hamburgerMenu = document.getElementById("hamburger-menu");
var navMenu = document.getElementById("nav-menu");
var nav = document.querySelector("nav");
var body = document.body;

if (hamburgerMenu && navMenu) {
  function toggleMenu() {
    var isActive = hamburgerMenu.classList.contains("active");

    if (isActive) {
      // Close menu
      hamburgerMenu.classList.remove("active");
      navMenu.classList.remove("active");
      nav.classList.remove("menu-open");
      body.classList.remove("menu-open");
    } else {
      // Open menu
      hamburgerMenu.classList.add("active");
      navMenu.classList.add("active");
      nav.classList.add("menu-open");
      body.classList.add("menu-open");
    }
  }

  function closeMenu() {
    hamburgerMenu.classList.remove("active");
    navMenu.classList.remove("active");
    nav.classList.remove("menu-open");
    body.classList.remove("menu-open");
  }

  hamburgerMenu.addEventListener("click", toggleMenu);

  // Close menu when clicking on a nav link (for better UX)
  var navLinks = navMenu.querySelectorAll("a");
  for (var i = 0; i < navLinks.length; i++) {
    navLinks[i].addEventListener("click", closeMenu);
  }

  // Close menu when clicking outside
  document.addEventListener("click", function (event) {
    var isClickInsideNav = navMenu.contains(event.target);
    var isClickOnHamburger = hamburgerMenu.contains(event.target);

    if (
      !isClickInsideNav &&
      !isClickOnHamburger &&
      navMenu.classList.contains("active")
    ) {
      closeMenu();
    }
  });
}
