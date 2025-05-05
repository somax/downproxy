package static

// IndexHTML 包含主页的 HTML 内容
const IndexHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>下载代理</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            max-width: 800px;
            margin: 20px auto;
            padding: 0 20px;
            background-color: #f5f5f5;
        }
        .container {
            background-color: white;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        h1 {
            color: #333;
            text-align: center;
            margin-bottom: 30px;
        }
        .input-group {
            display: flex;
            gap: 10px;
            margin-bottom: 20px;
        }
        input[type="url"] {
            flex: 1;
            padding: 10px;
            border: 1px solid #ddd;
            border-radius: 4px;
            font-size: 16px;
        }
        button {
            padding: 10px 20px;
            background-color: #007bff;
            color: white;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            font-size: 16px;
            transition: background-color 0.2s;
        }
        button:hover {
            background-color: #0056b3;
        }
        button:disabled {
            background-color: #ccc;
            cursor: not-allowed;
        }
        #status {
            margin-top: 20px;
            padding: 10px;
            border-radius: 4px;
            display: none;
        }
        .error {
            background-color: #fff3f3;
            color: #dc3545;
            border: 1px solid #dc3545;
        }
        .success {
            background-color: #f3fff3;
            color: #28a745;
            border: 1px solid #28a745;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>下载代理</h1>
        <div class="input-group">
            <input type="url" id="downloadUrl" placeholder="请输入下载链接" required>
            <button onclick="startDownload()" id="downloadBtn">开始下载</button>
        </div>
        <div id="status"></div>
    </div>

    <script>
        function startDownload() {
            const urlInput = document.getElementById('downloadUrl');
            const downloadBtn = document.getElementById('downloadBtn');
            const status = document.getElementById('status');
            const url = urlInput.value.trim();

            if (!url) {
                showStatus('请输入下载链接', 'error');
                return;
            }

            downloadBtn.disabled = true;
            const proxyUrl = '/download?url=' + encodeURIComponent(url);
            window.location.href = proxyUrl;
            showStatus('下载已开始...', 'success');
            
            setTimeout(() => {
                downloadBtn.disabled = false;
                status.style.display = 'none';
            }, 3000);
        }

        function showStatus(message, type) {
            const status = document.getElementById('status');
            status.textContent = message;
            status.className = type;
            status.style.display = 'block';
        }
    </script>
</body>
</html>`
