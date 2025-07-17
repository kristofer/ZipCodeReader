// Basic JavaScript for ZipCodeReader
// This will be expanded in later phases

document.addEventListener('DOMContentLoaded', function() {
    console.log('ZipCodeReader loaded successfully');
    
    // Add fade-in animation to main content
    const mainContent = document.querySelector('main');
    if (mainContent) {
        mainContent.classList.add('fade-in');
    }
    
    // Add hover effects to buttons
    const buttons = document.querySelectorAll('button');
    buttons.forEach(button => {
        button.classList.add('hover-scale');
    });
});

// Health check function for future use
async function checkHealth() {
    try {
        const response = await fetch('/health');
        const data = await response.json();
        console.log('Health check:', data);
        return data;
    } catch (error) {
        console.error('Health check failed:', error);
        return { status: 'error', message: 'Failed to check health' };
    }
}

// Utility function to show notifications (for future use)
function showNotification(message, type = 'info') {
    console.log(`${type.toUpperCase()}: ${message}`);
    // Will implement actual notification system in later phases
}
