document.addEventListener('DOMContentLoaded', () => {
    document.addEventListener('touchstart', (e) => {
        e.preventDefault();
        const touchEvent = e.touches[0];
        if (e.touches.length === 1) {
            console.log("Single touch detected: Simulating right-click.");
            const simulatedEvent = new MouseEvent('mousedown', {
                bubbles: true,
                cancelable: true,
                clientX: touchEvent.clientX,
                clientY: touchEvent.clientY,
                button: 2, // Right mouse button
            });
            touchEvent.target.dispatchEvent(simulatedEvent);
        } else if (e.touches.length > 1) {
            console.log("Multi-touch detected: Simulating left-click.");
            const simulatedEvent = new MouseEvent('mousedown', {
                bubbles: true,
                cancelable: true,
                clientX: touchEvent.clientX,
                clientY: touchEvent.clientY,
                button: 0, // Left mouse button
            });
            touchEvent.target.dispatchEvent(simulatedEvent);
        }
        console.log("Technical zero-move.");
        const simulatedEvent = new MouseEvent('mousemove', {
            bubbles: true,
            cancelable: true,
            clientX: touchEvent.clientX,
            clientY: touchEvent.clientY,
        });
        touchEvent.target.dispatchEvent(simulatedEvent);
    });

    document.addEventListener('touchmove', (e) => {
        e.preventDefault();
        console.log("Touch move detected.");
        const touchEvent = e.touches[0];
        const simulatedEvent = new MouseEvent('mousemove', {
            bubbles: true,
            cancelable: true,
            clientX: touchEvent.clientX,
            clientY: touchEvent.clientY,
        });
        touchEvent.target.dispatchEvent(simulatedEvent);
    });

    document.addEventListener('touchend', (e) => {
        e.preventDefault();
        console.log("Touch end detected.");
        const touchEvent = e.changedTouches[0];
        const simulatedEvent = new MouseEvent('mouseup', {
            bubbles: true,
            cancelable: true,
            clientX: touchEvent.clientX,
            clientY: touchEvent.clientY,
        });
        touchEvent.target.dispatchEvent(simulatedEvent);
    });
});
