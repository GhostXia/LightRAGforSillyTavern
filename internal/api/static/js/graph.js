// 知识图谱可视化脚本
let network = null;
let nodes = null;
let edges = null;
let allNodes = null;
let highlightActive = false;
let nodeColors = {};

// 页面加载完成后执行
document.addEventListener('DOMContentLoaded', function() {
    // 加载图谱列表
    loadGraphList();
    
    // 设置加载图谱按钮事件
    document.getElementById('loadGraphBtn').addEventListener('click', loadSelectedGraph);
    
    // 设置搜索按钮事件
    document.getElementById('searchBtn').addEventListener('click', searchNode);
    
    // 设置搜索框回车事件
    document.getElementById('searchNode').addEventListener('keypress', function(e) {
        if (e.key === 'Enter') {
            searchNode();
        }
    });
});

// 加载图谱列表
function loadGraphList() {
    fetch('/api/graphs')
        .then(response => response.json())
        .then(data => {
            const select = document.getElementById('graphSelect');
            // 清空现有选项（保留第一个默认选项）
            while (select.options.length > 1) {
                select.remove(1);
            }
            
            // 添加新选项
            if (data.graphs && data.graphs.length > 0) {
                data.graphs.forEach(graph => {
                    const option = document.createElement('option');
                    option.value = graph.id;
                    option.text = graph.name;
                    select.appendChild(option);
                });
            } else {
                // 如果没有图谱，添加提示选项
                const option = document.createElement('option');
                option.value = '';
                option.text = '没有可用的知识图谱';
                option.disabled = true;
                select.appendChild(option);
            }
        })
        .catch(error => {
            console.error('加载图谱列表失败:', error);
            showAlert('加载图谱列表失败: ' + error.message, 'danger');
        });
}

// 加载选中的图谱
function loadSelectedGraph() {
    const select = document.getElementById('graphSelect');
    const graphId = select.value;
    
    if (!graphId || graphId === '请选择知识图谱...') {
        showAlert('请选择一个知识图谱', 'warning');
        return;
    }
    
    fetch(`/api/graphs/${graphId}`)
        .then(response => response.json())
        .then(data => {
            if (data.nodes && data.edges) {
                drawGraph(data.nodes, data.edges);
            } else {
                showAlert('图谱数据格式不正确', 'danger');
            }
        })
        .catch(error => {
            console.error('加载图谱失败:', error);
            showAlert('加载图谱失败: ' + error.message, 'danger');
        });
}

// 绘制图谱
function drawGraph(nodeData, edgeData) {
    const container = document.getElementById('mynetwork');
    
    // 创建数据集
    nodes = new vis.DataSet(nodeData);
    edges = new vis.DataSet(edgeData);
    
    // 保存节点原始颜色
    allNodes = nodes.get({returnType: "Object"});
    for (let nodeId in allNodes) {
        nodeColors[nodeId] = allNodes[nodeId].color;
    }
    
    // 创建网络
    const data = {
        nodes: nodes,
        edges: edges
    };
    
    // 配置选项
    const options = {
        nodes: {
            shape: 'dot',
            size: 16,
            font: {
                size: 14,
                face: 'Microsoft YaHei'
            },
            borderWidth: 2
        },
        edges: {
            width: 1,
            font: {
                size: 12,
                align: 'middle',
                face: 'Microsoft YaHei'
            },
            arrows: {
                to: { enabled: true, scaleFactor: 0.5 }
            },
            smooth: { type: 'continuous' }
        },
        physics: {
            stabilization: false,
            barnesHut: {
                gravitationalConstant: -80000,
                springConstant: 0.001,
                springLength: 200
            }
        },
        interaction: {
            navigationButtons: true,
            keyboard: true
        }
    };
    
    // 创建网络
    network = new vis.Network(container, data, options);
    
    // 添加选择事件
    network.on("click", function(params) {
        if (params.nodes.length > 0) {
            neighbourhoodHighlight(params);
        } else {
            // 如果点击空白处，重置高亮
            resetHighlight();
        }
    });
}

// 搜索节点
function searchNode() {
    const searchText = document.getElementById('searchNode').value.toLowerCase();
    
    if (!searchText || !nodes) {
        return;
    }
    
    // 查找匹配的节点
    const matchedNodes = [];
    nodes.forEach(node => {
        if (node.label && node.label.toLowerCase().includes(searchText)) {
            matchedNodes.push(node.id);
        }
    });
    
    if (matchedNodes.length > 0) {
        // 选中第一个匹配的节点
        network.selectNodes(matchedNodes);
        network.focus(matchedNodes[0], {
            scale: 1.2,
            animation: true
        });
        neighbourhoodHighlight({nodes: matchedNodes});
    } else {
        showAlert('未找到匹配的节点', 'warning');
    }
}

// 节点高亮
function neighbourhoodHighlight(params) {
    // 如果选中了节点
    if (params.nodes.length > 0) {
        highlightActive = true;
        const selectedNode = params.nodes[0];
        const degrees = 2;
        
        // 将所有节点标记为难以阅读
        for (let nodeId in allNodes) {
            allNodes[nodeId].color = "rgba(200,200,200,0.5)";
            if (allNodes[nodeId].hiddenLabel === undefined) {
                allNodes[nodeId].hiddenLabel = allNodes[nodeId].label;
                allNodes[nodeId].label = undefined;
            }
        }
        
        // 获取连接的节点
        const connectedNodes = network.getConnectedNodes(selectedNode);
        let allConnectedNodes = [];
        
        // 获取二度节点
        for (let i = 1; i < degrees; i++) {
            for (let j = 0; j < connectedNodes.length; j++) {
                allConnectedNodes = allConnectedNodes.concat(
                    network.getConnectedNodes(connectedNodes[j])
                );
            }
        }
        
        // 二度节点使用不同颜色并恢复标签
        for (let i = 0; i < allConnectedNodes.length; i++) {
            allNodes[allConnectedNodes[i]].color = "rgba(150,150,150,0.75)";
            if (allNodes[allConnectedNodes[i]].hiddenLabel !== undefined) {
                allNodes[allConnectedNodes[i]].label = allNodes[allConnectedNodes[i]].hiddenLabel;
                allNodes[allConnectedNodes[i]].hiddenLabel = undefined;
            }
        }
        
        // 一度节点恢复原始颜色和标签
        for (let i = 0; i < connectedNodes.length; i++) {
            allNodes[connectedNodes[i]].color = nodeColors[connectedNodes[i]];
            if (allNodes[connectedNodes[i]].hiddenLabel !== undefined) {
                allNodes[connectedNodes[i]].label = allNodes[connectedNodes[i]].hiddenLabel;
                allNodes[connectedNodes[i]].hiddenLabel = undefined;
            }
        }
        
        // 选中的节点恢复原始颜色和标签
        allNodes[selectedNode].color = nodeColors[selectedNode];
        if (allNodes[selectedNode].hiddenLabel !== undefined) {
            allNodes[selectedNode].label = allNodes[selectedNode].hiddenLabel;
            allNodes[selectedNode].hiddenLabel = undefined;
        }
    }
    
    // 更新节点
    const updateArray = [];
    for (let nodeId in allNodes) {
        if (allNodes.hasOwnProperty(nodeId)) {
            updateArray.push(allNodes[nodeId]);
        }
    }
    nodes.update(updateArray);
}

// 重置高亮
function resetHighlight() {
    if (highlightActive === true) {
        // 重置所有节点
        for (let nodeId in allNodes) {
            allNodes[nodeId].color = nodeColors[nodeId];
            if (allNodes[nodeId].hiddenLabel !== undefined) {
                allNodes[nodeId].label = allNodes[nodeId].hiddenLabel;
                allNodes[nodeId].hiddenLabel = undefined;
            }
        }
        highlightActive = false;
        
        // 更新节点
        const updateArray = [];
        for (let nodeId in allNodes) {
            if (allNodes.hasOwnProperty(nodeId)) {
                updateArray.push(allNodes[nodeId]);
            }
        }
        nodes.update(updateArray);
    }
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
    const container = document.querySelector('.card-body');
    container.insertBefore(alertDiv, container.firstChild);
    
    // 3秒后自动关闭
    setTimeout(() => {
        alertDiv.classList.remove('show');
        setTimeout(() => alertDiv.remove(), 150);
    }, 3000);
}