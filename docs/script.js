// Classic Desktop JavaScript functionality
let activeWindow = null;
let zIndex = 1000;

// Window management
function openWindow(windowId) {
    const window = document.getElementById(windowId);
    if (window) {
        window.classList.add('active');
        window.style.zIndex = ++zIndex;
        centerWindow(window);
        activeWindow = window;
    }
}

function closeWindow(windowId) {
    const window = document.getElementById(windowId);
    if (window) {
        window.classList.remove('active');
    }
}

function minimizeWindow(windowId) {
    closeWindow(windowId);
}

function centerWindow(window) {
    const rect = window.getBoundingClientRect();
    window.style.left = (window.innerWidth - rect.width) / 2 + 'px';
    window.style.top = (window.innerHeight - rect.height) / 2 + 'px';
}

// Make windows draggable
function makeDraggable(element) {
    let pos1 = 0, pos2 = 0, pos3 = 0, pos4 = 0;
    const titleBar = element.querySelector('.title-bar');
    
    if (titleBar) {
        titleBar.onmousedown = dragMouseDown;
    }

    function dragMouseDown(e) {
        e = e || window.event;
        e.preventDefault();
        pos3 = e.clientX;
        pos4 = e.clientY;
        document.onmouseup = closeDragElement;
        document.onmousemove = elementDrag;
        
        // Bring to front
        element.style.zIndex = ++zIndex;
    }

    function elementDrag(e) {
        e = e || window.event;
        e.preventDefault();
        pos1 = pos3 - e.clientX;
        pos2 = pos4 - e.clientY;
        pos3 = e.clientX;
        pos4 = e.clientY;
        element.style.top = (element.offsetTop - pos2) + "px";
        element.style.left = (element.offsetLeft - pos1) + "px";
    }

    function closeDragElement() {
        document.onmouseup = null;
        document.onmousemove = null;
    }
}

// Update taskbar time
function updateTime() {
    const timeElement = document.getElementById('taskbar-time');
    if (timeElement) {
        const now = new Date();
        const hours = now.getHours();
        const minutes = now.getMinutes().toString().padStart(2, '0');
        const ampm = hours >= 12 ? 'PM' : 'AM';
        const displayHours = hours % 12 || 12;
        timeElement.textContent = `${displayHours}:${minutes} ${ampm}`;
    }
}

// Tab functionality
function switchTab(tabId, contentId) {
    // Hide all tab contents
    const tabContents = document.querySelectorAll('.tab-content');
    tabContents.forEach(content => {
        content.style.display = 'none';
    });
    
    // Remove active class from all tabs
    const tabs = document.querySelectorAll('.tab');
    tabs.forEach(tab => {
        tab.classList.remove('active');
    });
    
    // Show selected tab content
    const selectedContent = document.getElementById(contentId);
    if (selectedContent) {
        selectedContent.style.display = 'block';
    }
    
    // Add active class to selected tab
    const selectedTab = document.getElementById(tabId);
    if (selectedTab) {
        selectedTab.classList.add('active');
    }
}

// Copy to clipboard functionality
function copyToClipboard(text) {
    const textarea = document.createElement('textarea');
    textarea.value = text;
    document.body.appendChild(textarea);
    textarea.select();
    document.execCommand('copy');
    document.body.removeChild(textarea);
    
    // Show copied message
    alert('Copied to clipboard!');
}

// Initialize
document.addEventListener('DOMContentLoaded', function() {
    // Make all windows draggable
    const windows = document.querySelectorAll('.window');
    windows.forEach(window => {
        makeDraggable(window);
    });
    
    // Update time every second
    updateTime();
    setInterval(updateTime, 1000);
    
    // Initialize first tab if exists
    const firstTab = document.querySelector('.tab');
    if (firstTab) {
        firstTab.click();
    }
});