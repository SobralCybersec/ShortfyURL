const urlsTableBody = document.getElementById('urlsTableBody');
const totalUrls = document.getElementById('totalUrls');
const totalClicks = document.getElementById('totalClicks');

async function loadURLs() {
    try {
        const response = await fetch('/api/urls');
        
        if (!response.ok) {
            throw new Error('Erro ao carregar URLs');
        }

        const urls = await response.json();
        
        if (!urls || urls.length === 0) {
            urlsTableBody.innerHTML = '<tr><td colspan="5" class="loading">Nenhuma URL encontrada</td></tr>';
            return;
        }

        let clicks = 0;
        
        urlsTableBody.innerHTML = urls.map(url => {
            clicks += url.clicks;
            const date = new Date(url.created_at).toLocaleString('pt-BR');
            
            return `
                <tr>
                    <td><strong>${url.short_code}</strong></td>
                    <td class="url-cell" title="${url.original_url}">${url.original_url}</td>
                    <td>${url.clicks}</td>
                    <td>${date}</td>
                    <td>
                        <button class="action-btn" onclick="window.open('/${url.short_code}', '_blank')">Abrir</button>
                        <button class="action-btn" onclick="copyToClipboard('${window.location.origin}/${url.short_code}')">Copiar</button>
                    </td>
                </tr>
            `;
        }).join('');

        totalUrls.textContent = urls.length;
        totalClicks.textContent = clicks;

    } catch (err) {
        urlsTableBody.innerHTML = '<tr><td colspan="5" class="loading">Erro ao carregar dados</td></tr>';
        console.error(err);
    }
}

function copyToClipboard(text) {
    navigator.clipboard.writeText(text).then(() => {
        alert('Link copiado!');
    });
}

loadURLs();
setInterval(loadURLs, 5000);
