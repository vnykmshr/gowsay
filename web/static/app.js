// Theme management
function toggleTheme() {
    const html = document.documentElement;
    const currentTheme = html.getAttribute('data-theme');
    const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
    html.setAttribute('data-theme', newTheme);
    localStorage.setItem('theme', newTheme);
}

// Initialize theme from localStorage
function initTheme() {
    const savedTheme = localStorage.getItem('theme') || 'light';
    document.documentElement.setAttribute('data-theme', savedTheme);
}

// Load cows and moods from API
async function loadOptions() {
    try {
        const [cowsResp, moodsResp] = await Promise.all([
            fetch('/api/cows'),
            fetch('/api/moods')
        ]);

        const cowsData = await cowsResp.json();
        const moodsData = await moodsResp.json();

        const cowSelect = document.getElementById('cow');
        const moodSelect = document.getElementById('mood');

        // Populate cows
        cowSelect.innerHTML = '<option value="random">🎲 Random</option>';
        cowsData.cows.forEach(cow => {
            const option = document.createElement('option');
            option.value = cow;
            option.textContent = cow;
            cowSelect.appendChild(option);
        });

        // Populate moods
        moodSelect.innerHTML = '<option value="">Normal</option><option value="random">🎲 Random</option>';
        moodsData.moods.forEach(mood => {
            const option = document.createElement('option');
            option.value = mood;
            option.textContent = mood.charAt(0).toUpperCase() + mood.slice(1);
            moodSelect.appendChild(option);
        });
    } catch (error) {
        console.error('Failed to load options:', error);
    }
}

// Randomize selections
function randomize() {
    const cowSelect = document.getElementById('cow');
    const moodSelect = document.getElementById('mood');
    const actionSelect = document.getElementById('action');

    cowSelect.value = 'random';
    moodSelect.value = 'random';
    actionSelect.value = Math.random() > 0.5 ? 'say' : 'think';

    const messages = [
        "Moo!",
        "Have a nice day!",
        "Time for coffee ☕",
        "Hello, World!",
        "Keep calm and moo on",
        "Life is better with cows",
        "Udderly awesome!",
        "Got milk?",
        "Moo-velous day ahead!",
        "Don't have a cow, man!"
    ];

    document.getElementById('text').value = messages[Math.floor(Math.random() * messages.length)];
}

// Handle form submission
document.getElementById('cowform').addEventListener('submit', async (e) => {
    e.preventDefault();

    const submitBtn = e.target.querySelector('button[type="submit"]');
    const outputCard = document.getElementById('output-card');
    const output = document.getElementById('output');

    const formData = {
        text: document.getElementById('text').value,
        cow: document.getElementById('cow').value,
        mood: document.getElementById('mood').value,
        action: document.getElementById('action').value
    };

    // Show loading state
    submitBtn.classList.add('loading');
    submitBtn.disabled = true;

    try {
        const response = await fetch('/api/moo', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(formData)
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Failed to generate cowsay');
        }

        const data = await response.json();

        // Strip markdown code blocks if present
        let outputText = data.output;
        outputText = outputText.replace(/^```\n/, '').replace(/\n```\n$/, '');

        output.textContent = outputText;
        outputCard.style.display = 'block';

        // Smooth scroll to output
        outputCard.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
    } catch (error) {
        output.textContent = `Error: ${error.message}`;
        outputCard.style.display = 'block';
    } finally {
        submitBtn.classList.remove('loading');
        submitBtn.disabled = false;
    }
});

// Copy to clipboard
async function copyToClipboard() {
    const output = document.getElementById('output');
    const btn = event.currentTarget;

    try {
        await navigator.clipboard.writeText(output.textContent);

        // Show success feedback
        btn.classList.add('copied');
        setTimeout(() => {
            btn.classList.remove('copied');
        }, 2000);
    } catch (error) {
        console.error('Failed to copy:', error);
        // Fallback for older browsers
        const textArea = document.createElement('textarea');
        textArea.value = output.textContent;
        textArea.style.position = 'fixed';
        textArea.style.opacity = '0';
        document.body.appendChild(textArea);
        textArea.select();
        try {
            document.execCommand('copy');
            btn.classList.add('copied');
            setTimeout(() => {
                btn.classList.remove('copied');
            }, 2000);
        } catch (err) {
            console.error('Fallback copy failed:', err);
        }
        document.body.removeChild(textArea);
    }
}

// Initialize on page load
document.addEventListener('DOMContentLoaded', () => {
    initTheme();
    loadOptions();

    // Set default message
    const defaultMessages = [
        "Hello, World!",
        "Welcome to gowsay!",
        "Moo!",
    ];
    document.getElementById('text').value = defaultMessages[Math.floor(Math.random() * defaultMessages.length)];
});
