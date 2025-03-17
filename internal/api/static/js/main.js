// 页面加载完成后执行
document.addEventListener('DOMContentLoaded', function() {
    // 加载配置信息
    loadConfig();
    
    // 设置表单提交事件
    document.getElementById('settingsForm').addEventListener('submit', function(e) {
        e.preventDefault();
        saveConfig();
    });
    
    // 设置发送消息按钮事件
    document.getElementById('sendBtn').addEventListener('click', sendMessage);
    
    // 设置消息输入框回车事件
    document.getElementById('messageInput').addEventListener('keypress', function(e) {
        if (e.key === 'Enter') {
            sendMessage();
        }
    });
    
    // 设置上传按钮事件
    document.getElementById('uploadBtn').addEventListener('click', uploadFiles);
    
    // 加载文档列表
    loadDocuments();
});

// 加载配置信息
function loadConfig() {
    fetch('/web/config')
        .then(response => response.json())
        .then(data => {
            // 填充RAG设置
            document.getElementById('workingDir').value = data.rag.workingDir || '';
            document.getElementById('inputDir').value = data.rag.inputDir || '';
            document.getElementById('maxTokens').value = data.rag.maxTokens || 0;
            document.getElementById('maxEmbedTokens').value = data.rag.maxEmbedTokens || 0;
            document.getElementById('chunkSize').value = data.rag.chunkSize || 0;
            document.getElementById('chunkOverlap').value = data.rag.chunkOverlap || 0;
            
            // 填充LLM设置
            document.getElementById('llmProvider').value = data.llm.provider || 'openai';
            document.getElementById('llmModel').value = data.llm.model || '';
            if (document.getElementById('llmBaseUrl')) {
                document.getElementById('llmBaseUrl').value = data.llm.baseUrl || '';
            }
            
            // 填充嵌入模型设置
            document.getElementById('embeddingProvider').value = data.embedding.provider || 'openai';
            document.getElementById('embeddingModel').value = data.embedding.model || '';
            if (document.getElementById('embeddingBaseUrl')) {
                document.getElementById('embeddingBaseUrl').value = data.embedding.baseUrl || '';
            }
            
            // 填充存储设置
            document.getElementById('vectorType').value = data.storage.vectorType || 'memory';
            document.getElementById('kvType').value = data.storage.kvType || 'memory';
            document.getElementById('graphType').value = data.storage.graphType || 'memory';
        })
        .catch(error => {
            console.error('加载配置失败:', error);
            showAlert('加载配置失败: ' + error.message, 'danger');
        });
}

// 保存配置信息
function saveConfig() {
    const formData = {
        rag: {
            workingDir: document.getElementById('workingDir').value,
            inputDir: document.getElementById('inputDir').value,
            maxTokens: parseInt(document.getElementById('maxTokens').value),
            maxEmbedTokens: parseInt(document.getElementById('maxEmbedTokens').value),
            chunkSize: parseInt(document.getElementById('chunkSize').value),
            chunkOverlap: parseInt(document.getElementById('chunkOverlap').value)
        },
        llm: {
            provider: document.getElementById('llmProvider').value,
            model: document.getElementById('llmModel').value
        },
        embedding: {
            provider: document.getElementById('embeddingProvider').value,
            model: document.getElementById('embeddingModel').value
        },
        storage: {
            vectorType: document.getElementById('vectorType').value,
            kvType: document.getElementById('kvType').value,
            graphType: document.getElementById('graphType').value
        },
        server: {}
    };
    
    // 添加可选的API密钥和URL
    if (document.getElementById('llmApiKey') && document.getElementById('llmApiKey').value) {
        formData.llm.apiKey = document.getElementById('llmApiKey').value;
    }
    
    if (document.getElementById('llmBaseUrl') && document.getElementById('llmBaseUrl').value) {
        formData.llm.baseUrl = document.getElementById('llmBaseUrl').value;
    }
    
    if (document.getElementById('embeddingApiKey') && document.getElementById('embeddingApiKey').value) {
        formData.embedding.apiKey = document.getElementById('embeddingApiKey').value;
    }
    
    if (document.getElementById('embeddingBaseUrl') && document.getElementById('embeddingBaseUrl').value) {
        formData.embedding.baseUrl = document.getElementById('embeddingBaseUrl').value;
    }
    
    if (document.getElementById('serverApiKey') && document.getElementById('serverApiKey').value) {
        formData.server.apiKey = document.getElementById('serverApiKey').value;
    }
    
    fetch('/web/config', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(formData)
    })
    .then(response => response.json())
    .then(data => {
        showAlert('配置已保存', 'success');
    })
    .catch(error => {
        console.error('保存配置失败:', error);
        showAlert('保存配置失败: ' + error.message, 'danger');
    });
}

// 发送消息
function sendMessage() {
    const messageInput = document.getElementById('messageInput');
    const message = messageInput.value.trim();
    
    if (message === '') return;
    
    // 添加用户消息到聊天窗口
    addMessage('user', message);
    
    // 清空输入框
    messageInput.value = '';
    
    // 显示加载中
    const loadingId = showLoading();
    
    // 发送请求到服务器
    fetch('/api/query', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ query: message })
    })
    .then(response => response.json())
    .then(data => {
        // 隐藏加载中
        hideLoading(loadingId);
        
        // 添加助手回复到聊天窗口
        addMessage('assistant', data.response || data.message);
    })
    .catch(error => {
        // 隐藏加载中
        hideLoading(loadingId);
        
        console.error('查询失败:', error);
        addMessage('system', '查询失败: ' + error.message);
    });
}

// 添加消息到聊天窗口
function addMessage(role, content) {
    const chatMessages = document.getElementById('chatMessages');
    const messageDiv = document.createElement('div');
    
    messageDiv.className = role + '-message';
    messageDiv.textContent = content;
    
    chatMessages.appendChild(messageDiv);
    
    // 滚动到底部
    chatMessages.scrollTop = chatMessages.scrollHeight;
}

// 显示加载中
function showLoading() {
    const chatMessages = document.getElementById('chatMessages');
    const loadingDiv = document.createElement('div');
    const loadingId = 'loading-' + Date.now();
    
    loadingDiv.className = 'system-message';
    loadingDiv.id = loadingId;
    loadingDiv.textContent = '正在思考...';
    
    chatMessages.appendChild(loadingDiv);
    chatMessages.scrollTop = chatMessages.scrollHeight;
    
    return loadingId;
}

// 隐藏加载中
function hideLoading(loadingId) {
    const loadingDiv = document.getElementById(loadingId);
    if (loadingDiv) {
        loadingDiv.remove();
    }
}

// 上传文件
function uploadFiles() {
    const fileInput = document.getElementById('fileUpload');
    const files = fileInput.files;
    
    if (files.length === 0) {
        showAlert('请选择要上传的文件', 'warning');
        return;
    }
    
    const formData = new FormData();
    for (let i = 0; i < files.length; i++) {
        formData.append('files', files[i]);
    }
    
    // 显示上传中提示
    showAlert('文件上传中...', 'info');
    
    fetch('/api/upload', {
        method: 'POST',
        body: formData
    })
    .then(response => response.json())
    .then(data => {
        showAlert('文件上传成功', 'success');
        fileInput.value = ''; // 清空文件选择
        loadDocuments(); // 重新加载文档列表
    })
    .catch(error => {
        console.error('上传失败:', error);
        showAlert('上传失败: ' + error.message, 'danger');
    });
}

// 加载文档列表
function loadDocuments() {
    fetch('/api/documents')
        .then(response => response.json())
        .then(data => {
            const tableBody = document.querySelector('#documentTable tbody');
            tableBody.innerHTML = ''; // 清空表格
            
            if (data.documents && data.documents.length > 0) {
                data.documents.forEach(doc => {
                    const row = document.createElement('tr');
                    
                    // 文件名
                    const nameCell = document.createElement('td');
                    nameCell.textContent = doc.name;
                    row.appendChild(nameCell);
                    
                    // 文件大小
                    const sizeCell = document.createElement('td');
                    sizeCell.textContent = formatFileSize(doc.size);
                    row.appendChild(sizeCell);
                    
                    // 上传时间
                    const timeCell = document.createElement('td');
                    timeCell.textContent = new Date(doc.uploadTime).toLocaleString();
                    row.appendChild(timeCell);
                    
                    // 操作按钮
                    const actionCell = document.createElement('td');
                    const deleteBtn = document.createElement('button');
                    deleteBtn.className = 'btn btn-sm btn-danger';
                    deleteBtn.textContent = '删除';
                    deleteBtn.onclick = function() {
                        deleteDocument(doc.id);
                    };
                    actionCell.appendChild(deleteBtn);
                    row.appendChild(actionCell);
                    
                    tableBody.appendChild(row);
                });
            } else {
                // 没有文档时显示提示
                const row = document.createElement('tr');
                const cell = document.createElement('td');
                cell.colSpan = 4;
                cell.textContent = '暂无文档';
                cell.className = 'text-center';
                row.appendChild(cell);
                tableBody.appendChild(row);
            }
        })
        .catch(error => {
            console.error('加载文档列表失败:', error);
        });
}

// 删除文档
function deleteDocument(id) {
    if (!confirm('确定要删除此文档吗？')) {
        return;
    }
    
    fetch(`/api/documents/${id}`, {
        method: 'DELETE'
    })
    .then(response => response.json())
    .then(data => {
        showAlert('文档已删除', 'success');
        loadDocuments(); // 重新加载文档列表
    })
    .catch(error => {
        console.error('删除文档失败:', error);
        showAlert('删除文档失败: ' + error.message, 'danger');
    });
}

// 格式化文件大小
function formatFileSize(bytes) {
    if (bytes === 0) return '0 Bytes';
    
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

// 显示提示信息
function showAlert(message, type) {
    // 创建提示元素
    const alertDiv = document.createElement('div');
    alertDiv.className = `alert alert-${type} alert-dismissible fade show`;
    alertDiv.role = 'alert';
    alertDiv.innerHTML = `
        ${message}
        <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
    `;
    
    // 添加到页面
    const container = document.querySelector('.container');
    container.insertBefore(alertDiv, container.firstChild);
    
    // 3秒后自动关闭
    setTimeout(() => {
        alertDiv.classList.remove('show')