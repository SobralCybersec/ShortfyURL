const urlInput = document.getElementById('urlInput');
const shortenBtn = document.getElementById('shortenBtn');
const result = document.getElementById('result');
const error = document.getElementById('error');
const shortUrlLink = document.getElementById('shortUrlLink');
const copyBtn = document.getElementById('copyBtn');
const statsLink = document.getElementById('statsLink');

shortenBtn.addEventListener('click', async () => {
    const url = urlInput.value.trim();
    
    if (!url) {
        showError('Por favor, insira uma URL válida');
        return;
    }

    try {
        shortenBtn.disabled = true;
        shortenBtn.textContent = 'Encurtando...';

        const response = await fetch('/api/shorten', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ url })
        });

        if (!response.ok) {
            throw new Error('Erro ao encurtar URL');
        }

        const data = await response.json();
        
        shortUrlLink.href = data.short_url;
        shortUrlLink.textContent = data.short_url;
        statsLink.href = `/api/stats/${data.short_code}`;
        
        result.classList.remove('hidden');
        error.classList.add('hidden');
        urlInput.value = '';

    } catch (err) {
        showError('Erro ao encurtar URL. Tente novamente.');
    } finally {
        shortenBtn.disabled = false;
        shortenBtn.textContent = 'Encurtar';
    }
});

copyBtn.addEventListener('click', () => {
    navigator.clipboard.writeText(shortUrlLink.href).then(() => {
        const originalText = copyBtn.textContent;
        copyBtn.textContent = '✓ Copiado!';
        
        setTimeout(() => {
            copyBtn.textContent = originalText;
        }, 2000);
    });
});

urlInput.addEventListener('keypress', (e) => {
    if (e.key === 'Enter') {
        shortenBtn.click();
    }
});

function showError(message) {
    error.textContent = message;
    error.classList.remove('hidden');
    result.classList.add('hidden');
}
