// Beanstalkd Web UI - JavaScript

const API_BASE = '/api';
let currentTab = 'tubes';
let refreshInterval = null;

// 初始化
document.addEventListener('DOMContentLoaded', function() {
    loadInitialData();
    startAutoRefresh();
});

// 启动自动刷新
function startAutoRefresh() {
    refreshInterval = setInterval(() => {
        refreshAll();
    }, 10000); // 每10秒刷新一次
}

// 刷新所有数据
async function refreshAll() {
    await loadServerStats();
    await loadTubes();
    if (currentTab === 'jobs') {
        const tube = document.getElementById('jobsTubeSelect').value;
        if (tube) {
            await loadJobsForTube();
        }
    }
}

// 加载初始数据
async function loadInitialData() {
    await loadServerStats();
    await loadTubes();
    await populateTubeSelects();
}

// 切换标签
function switchTab(tabName) {
    currentTab = tabName;
    
    // 更新标签按钮样式
    document.querySelectorAll('.tab-button').forEach(btn => {
        btn.classList.remove('active', 'border-primary', 'text-primary');
        btn.classList.add('border-transparent', 'text-gray-500', 'hover:text-gray-700', 'hover:border-gray-300');
    });
    
    const activeBtn = document.getElementById(`tab-${tabName}`);
    activeBtn.classList.add('active', 'border-primary', 'text-primary');
    activeBtn.classList.remove('border-transparent', 'text-gray-500', 'hover:text-gray-700', 'hover:border-gray-300');
    
    // 显示对应内容
    document.querySelectorAll('.tab-content').forEach(content => {
        content.classList.add('hidden');
    });
    document.getElementById(`content-${tabName}`).classList.remove('hidden');
    
    // 加载对应数据
    if (tabName === 'stats') {
        loadServerStats();
    } else if (tabName === 'tubes') {
        loadTubes();
    }
}

// 加载服务器统计
async function loadServerStats() {
    try {
        const response = await fetch(`${API_BASE}/stats`);
        const data = await response.json();
        
        if (data.error) {
            showToast('错误', data.error, 'error');
            return;
        }
        
        // 更新顶部卡片
        document.getElementById('readyJobs').textContent = data.stats['current-jobs-ready'] || '0';
        document.getElementById('reservedJobs').textContent = data.stats['current-jobs-reserved'] || '0';
        document.getElementById('totalTubes').textContent = data.stats['current-tubes'] || '0';
        document.getElementById('totalJobs').textContent = data.stats['total-jobs'] || '0';
        
        // 更新详细统计
        const statsContainer = document.getElementById('serverStats');
        if (statsContainer && !statsContainer.querySelector('.stat-item')) {
            statsContainer.innerHTML = '';
            
            const sortedKeys = Object.keys(data.stats).sort();
            sortedKeys.forEach(key => {
                const value = data.stats[key];
                const statCard = document.createElement('div');
                statCard.className = 'stat-item bg-gradient-to-br from-gray-50 to-gray-100 rounded-lg p-4 border border-gray-200 hover:shadow-md transition-shadow';
                statCard.innerHTML = `
                    <p class="text-xs text-gray-500 font-medium uppercase tracking-wider">${key}</p>
                    <p class="text-lg font-bold text-gray-800 mt-1">${value}</p>
                `;
                statsContainer.appendChild(statCard);
            });
        }
    } catch (error) {
        console.error('加载服务器统计失败:', error);
        showToast('错误', '加载服务器统计失败', 'error');
    }
}

// 加载 Tubes 列表
async function loadTubes() {
    try {
        const response = await fetch(`${API_BASE}/tubes`);
        const data = await response.json();
        
        if (data.error) {
            showToast('错误', data.error, 'error');
            return;
        }
        
        const tubesContainer = document.getElementById('tubesList');
        tubesContainer.innerHTML = '';
        
        if (!data.tubes || data.tubes.length === 0) {
            tubesContainer.innerHTML = '<p class="text-center text-gray-500 py-8">暂无 Tubes</p>';
            return;
        }
        
        for (const tubeName of data.tubes) {
            const tubeCard = await createTubeCard(tubeName);
            tubesContainer.appendChild(tubeCard);
        }
        
        // 同时更新下拉选择框
        await populateTubeSelects();
    } catch (error) {
        console.error('加载 Tubes 失败:', error);
        showToast('错误', '加载 Tubes 失败', 'error');
    }
}

// 创建 Tube 卡片
async function createTubeCard(tubeName) {
    const card = document.createElement('div');
    card.className = 'bg-gradient-to-r from-white to-gray-50 rounded-lg p-6 border border-gray-200 hover:shadow-lg transition-all duration-200';
    
    try {
        const response = await fetch(`${API_BASE}/tubes/${tubeName}/stats`);
        const data = await response.json();
        
        if (data.error) {
            card.innerHTML = `
                <div class="flex items-center justify-between">
                    <h3 class="text-lg font-bold text-gray-800">${tubeName}</h3>
                    <span class="text-sm text-red-500">加载失败</span>
                </div>
            `;
            return card;
        }
        
        const stats = data.stats;
        card.innerHTML = `
            <div class="flex items-center justify-between mb-4">
                <h3 class="text-lg font-bold text-gray-800 flex items-center">
                    <i class="fas fa-layer-group text-primary mr-2"></i>
                    ${tubeName}
                </h3>
                <button onclick="viewTubeDetails('${tubeName}')" class="px-3 py-1 text-sm bg-blue-100 text-blue-700 rounded-lg hover:bg-blue-200 transition-colors">
                    <i class="fas fa-eye mr-1"></i>详情
                </button>
            </div>
            <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
                <div class="text-center">
                    <p class="text-2xl font-bold text-blue-600">${stats['current-jobs-ready'] || 0}</p>
                    <p class="text-xs text-gray-500 mt-1">就绪</p>
                </div>
                <div class="text-center">
                    <p class="text-2xl font-bold text-yellow-600">${stats['current-jobs-reserved'] || 0}</p>
                    <p class="text-xs text-gray-500 mt-1">预留</p>
                </div>
                <div class="text-center">
                    <p class="text-2xl font-bold text-red-600">${stats['current-jobs-buried'] || 0}</p>
                    <p class="text-xs text-gray-500 mt-1">埋葬</p>
                </div>
                <div class="text-center">
                    <p class="text-2xl font-bold text-purple-600">${stats['current-jobs-delayed'] || 0}</p>
                    <p class="text-xs text-gray-500 mt-1">延迟</p>
                </div>
            </div>
            <div class="mt-4 pt-4 border-t border-gray-200 flex justify-between text-sm text-gray-600">
                <span>总任务: <strong>${stats['total-jobs'] || 0}</strong></span>
                <span>监视者: <strong>${stats['current-watching'] || 0}</strong></span>
            </div>
        `;
    } catch (error) {
        card.innerHTML = `
            <div class="flex items-center justify-between">
                <h3 class="text-lg font-bold text-gray-800">${tubeName}</h3>
                <span class="text-sm text-red-500">加载失败</span>
            </div>
        `;
    }
    
    return card;
}

// 填充 Tube 下拉框
async function populateTubeSelects() {
    try {
        const response = await fetch(`${API_BASE}/tubes`);
        const data = await response.json();
        
        if (data.tubes) {
            const selects = ['jobsTubeSelect'];
            selects.forEach(selectId => {
                const select = document.getElementById(selectId);
                if (select) {
                    const currentValue = select.value;
                    select.innerHTML = '<option value="">选择一个 Tube...</option>';
                    data.tubes.forEach(tube => {
                        const option = document.createElement('option');
                        option.value = tube;
                        option.textContent = tube;
                        select.appendChild(option);
                    });
                    if (currentValue) {
                        select.value = currentValue;
                    }
                }
            });
        }
    } catch (error) {
        console.error('加载 Tube 列表失败:', error);
    }
}

// 加载指定 Tube 的任务
async function loadJobsForTube() {
    const tube = document.getElementById('jobsTubeSelect').value;
    if (!tube) return;
    
    const jobsContainer = document.getElementById('jobsList');
    jobsContainer.innerHTML = '<div class="text-center py-8 text-gray-500"><i class="fas fa-spinner fa-spin text-3xl mb-2"></i><p>加载中...</p></div>';
    
    try {
        // 获取 tube 统计信息
        const response = await fetch(`${API_BASE}/tubes/${tube}/stats`);
        const data = await response.json();
        
        if (data.error) {
            jobsContainer.innerHTML = `<p class="text-center text-red-500 py-8">${data.error}</p>`;
            return;
        }
        
        const stats = data.stats;
        jobsContainer.innerHTML = '';
        
        // 显示统计信息
        const statsCard = document.createElement('div');
        statsCard.className = 'bg-gradient-to-r from-blue-50 to-purple-50 rounded-lg p-6 mb-6 border border-blue-200';
        statsCard.innerHTML = `
            <h3 class="text-lg font-bold text-gray-800 mb-4">Tube: ${tube}</h3>
            <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
                <div class="text-center">
                    <p class="text-3xl font-bold text-blue-600">${stats['current-jobs-ready'] || 0}</p>
                    <p class="text-sm text-gray-600 mt-1">就绪任务</p>
                </div>
                <div class="text-center">
                    <p class="text-3xl font-bold text-yellow-600">${stats['current-jobs-reserved'] || 0}</p>
                    <p class="text-sm text-gray-600 mt-1">预留任务</p>
                </div>
                <div class="text-center">
                    <p class="text-3xl font-bold text-red-600">${stats['current-jobs-buried'] || 0}</p>
                    <p class="text-sm text-gray-600 mt-1">埋葬任务</p>
                </div>
                <div class="text-center">
                    <p class="text-3xl font-bold text-purple-600">${stats['current-jobs-delayed'] || 0}</p>
                    <p class="text-sm text-gray-600 mt-1">延迟任务</p>
                </div>
            </div>
        `;
        jobsContainer.appendChild(statsCard);
        
        // 提示信息
        const infoCard = document.createElement('div');
        infoCard.className = 'bg-blue-50 border border-blue-200 rounded-lg p-4';
        infoCard.innerHTML = `
            <div class="flex items-start">
                <i class="fas fa-info-circle text-blue-500 text-xl mr-3 mt-1"></i>
                <div>
                    <p class="text-sm text-gray-700"><strong>提示：</strong>使用"操作中心"标签来预留和查看任务详情</p>
                    <p class="text-sm text-gray-600 mt-1">总任务数: ${stats['total-jobs'] || 0}</p>
                </div>
            </div>
        `;
        jobsContainer.appendChild(infoCard);
        
    } catch (error) {
        console.error('加载任务失败:', error);
        jobsContainer.innerHTML = `<p class="text-center text-red-500 py-8">加载失败: ${error.message}</p>`;
    }
}

// 插入任务
async function putJob() {
    const tube = document.getElementById('putTube').value;
    const data = document.getElementById('putData').value;
    const priority = parseInt(document.getElementById('putPriority').value);
    const delay = parseInt(document.getElementById('putDelay').value);
    
    if (!data) {
        showToast('错误', '请输入任务数据', 'error');
        return;
    }
    
    try {
        const response = await fetch(`${API_BASE}/put`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ tube, data, priority, delay })
        });
        
        const result = await response.json();
        
        if (result.error) {
            showToast('错误', result.error, 'error');
        } else {
            showToast('成功', `任务已插入，ID: ${result.job_id}`, 'success');
            document.getElementById('putData').value = '';
            refreshAll();
        }
    } catch (error) {
        showToast('错误', '插入任务失败: ' + error.message, 'error');
    }
}

// 预留任务
async function reserveJob() {
    const tube = document.getElementById('reserveTube').value;
    const timeout = parseInt(document.getElementById('reserveTimeout').value);
    
    try {
        const response = await fetch(`${API_BASE}/reserve`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ tube, timeout })
        });
        
        const result = await response.json();
        
        if (result.error) {
            showToast('错误', result.error, 'error');
            document.getElementById('reserveResult').classList.add('hidden');
        } else {
            showToast('成功', `任务已预留，ID: ${result.job_id}`, 'success');
            document.getElementById('reserveResultContent').textContent = JSON.stringify({
                job_id: result.job_id,
                data: result.data
            }, null, 2);
            document.getElementById('reserveResult').classList.remove('hidden');
            refreshAll();
        }
    } catch (error) {
        showToast('错误', '预留任务失败: ' + error.message, 'error');
    }
}

// 删除任务
async function deleteJob() {
    const jobId = document.getElementById('deleteJobId').value;
    
    if (!jobId) {
        showToast('错误', '请输入任务 ID', 'error');
        return;
    }
    
    if (!confirm(`确定要删除任务 ${jobId} 吗？`)) {
        return;
    }
    
    try {
        const response = await fetch(`${API_BASE}/delete`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ job_id: parseInt(jobId) })
        });
        
        const result = await response.json();
        
        if (result.error) {
            showToast('错误', result.error, 'error');
        } else {
            showToast('成功', `任务 ${jobId} 已删除`, 'success');
            document.getElementById('deleteJobId').value = '';
            refreshAll();
        }
    } catch (error) {
        showToast('错误', '删除任务失败: ' + error.message, 'error');
    }
}

// 踢出任务
async function kickJobs() {
    const tube = document.getElementById('kickTube').value;
    const bound = parseInt(document.getElementById('kickBound').value);
    
    try {
        const response = await fetch(`${API_BASE}/kick`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ tube, bound })
        });
        
        const result = await response.json();
        
        if (result.error) {
            showToast('错误', result.error, 'error');
        } else {
            showToast('成功', `已踢出 ${result.kicked} 个任务`, 'success');
            refreshAll();
        }
    } catch (error) {
        showToast('错误', '踢出任务失败: ' + error.message, 'error');
    }
}

// 查看 Tube 详情
function viewTubeDetails(tubeName) {
    document.getElementById('jobsTubeSelect').value = tubeName;
    switchTab('jobs');
    loadJobsForTube();
}

// 显示 Toast 通知
function showToast(title, message, type = 'info') {
    const toast = document.getElementById('toast');
    const icon = document.getElementById('toastIcon');
    const titleEl = document.getElementById('toastTitle');
    const messageEl = document.getElementById('toastMessage');
    const toastDiv = toast.querySelector('div');
    
    // 设置图标和颜色
    const config = {
        success: { icon: 'fas fa-check-circle text-green-500', border: 'border-green-500' },
        error: { icon: 'fas fa-exclamation-circle text-red-500', border: 'border-red-500' },
        info: { icon: 'fas fa-info-circle text-blue-500', border: 'border-blue-500' }
    };
    
    const { icon: iconClass, border } = config[type] || config.info;
    icon.className = iconClass + ' text-2xl';
    toastDiv.className = `bg-white rounded-lg shadow-2xl p-4 ${border} border-l-4 min-w-[300px] transform transition-all duration-300`;
    
    titleEl.textContent = title;
    messageEl.textContent = message;
    
    // 显示Toast
    toast.classList.remove('hidden');
    
    // 3秒后自动隐藏
    setTimeout(() => {
        toast.classList.add('hidden');
    }, 3000);
}

// 添加初始标签样式
document.addEventListener('DOMContentLoaded', function() {
    const style = document.createElement('style');
    style.textContent = `
        .tab-button {
            transition: all 0.3s ease;
        }
        .tab-button.active {
            color: #3b82f6;
            border-color: #3b82f6;
        }
        .tab-button:not(.active) {
            color: #6b7280;
            border-color: transparent;
        }
        .tab-button:not(.active):hover {
            color: #374151;
            border-color: #d1d5db;
        }
    `;
    document.head.appendChild(style);
});
