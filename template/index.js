/**
 * RECON TOOL CMS - FINAL VERSION
 * Features: Config Management, Queue Monitor, Full SIM Management (CRUD, Bulk, Import/Export)
 */

const STYLES = `
:root {
  --bg-body: #020617; --bg-panel: #0f172a; --bg-input: #1e293b;
  --border-color: #334155; --text-main: #e2e8f0; --text-muted: #94a3b8;
  --primary: #3b82f6; --primary-hover: #2563eb; --danger: #ef4444;
  --success: #10b981; --warning: #facc15;
}
* { box-sizing: border-box; }
body { font-family: 'Inter', system-ui, sans-serif; color: var(--text-main); background: var(--bg-body); margin: 0; font-size: 14px; }

/* LAYOUT */
.app-container { display: flex; height: 100vh; }
.sidebar { width: 240px; background: var(--bg-panel); border-right: 1px solid var(--border-color); flex-shrink: 0; display: flex; flex-direction: column; }
.sidebar-header { padding: 24px; font-weight: 800; font-size: 18px; color: var(--primary); border-bottom: 1px solid var(--border-color); letter-spacing: -0.5px; }
.nav-item { padding: 14px 24px; cursor: pointer; color: var(--text-muted); transition: 0.2s; display: flex; align-items: center; gap: 10px; font-weight: 500; }
.nav-item:hover { background: rgba(255,255,255,0.03); color: var(--text-main); }
.nav-item.active { background: rgba(59, 130, 246, 0.1); color: var(--primary); border-right: 3px solid var(--primary); }

.main-content { flex: 1; display: flex; flex-direction: column; overflow: hidden; }
.view-area { flex: 1; padding: 30px; overflow-y: auto; }

/* COMPONENTS */
.panel { background: var(--bg-panel); border: 1px solid var(--border-color); border-radius: 12px; padding: 24px; margin-bottom: 24px; }
.grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(350px, 1fr)); gap: 24px; }

/* FORMS */
.form-group { margin-bottom: 16px; }
label { display: block; color: var(--text-muted); margin-bottom: 8px; font-size: 11px; font-weight: 700; text-transform: uppercase; letter-spacing: 0.5px; }
input, select { 
    width: 100%; background: var(--bg-input); border: 1px solid var(--border-color); 
    color: var(--text-main); padding: 10px 12px; border-radius: 8px; outline: none; font-size: 14px; transition: 0.2s;
}
input:focus, select:focus { border-color: var(--primary); box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.2); }
input:disabled { opacity: 0.6; cursor: not-allowed; }

/* BUTTONS */
button { 
    padding: 10px 16px; border-radius: 8px; cursor: pointer; border: none; font-weight: 600; 
    transition: 0.2s; background: var(--border-color); color: var(--text-main); font-size: 13px;
    display: inline-flex; align-items: center; justify-content: center; gap: 6px;
}
button:hover { filter: brightness(1.1); transform: translateY(-1px); }
button:active { transform: translateY(0); }
button.primary { background: var(--primary); color: white; }
button.danger { background: rgba(239, 68, 68, 0.15); color: var(--danger); }
button.danger:hover { background: var(--danger); color: white; }
button.success { background: rgba(16, 185, 129, 0.15); color: var(--success); }

/* TABLES */
table { width: 100%; border-collapse: collapse; }
th { text-align: left; padding: 12px; color: var(--text-muted); font-size: 11px; border-bottom: 1px solid var(--border-color); text-transform: uppercase; }
td { padding: 12px; border-bottom: 1px solid var(--border-color); vertical-align: middle; }
tr:hover td { background: rgba(255,255,255,0.02); }

.badge { padding: 4px 8px; border-radius: 6px; font-size: 11px; font-weight: 700; }
/* STATUS BADGES UPDATED */
.status-1 { background: rgba(148, 163, 184, 0.2); color: #94a3b8; border: 1px solid rgba(148, 163, 184, 0.3); } /* Normal - Xám */
.status-2 { background: rgba(16, 185, 129, 0.2); color: #34d399; border: 1px solid rgba(16, 185, 129, 0.3); } /* Has Plan - Xanh lá (Có gói) */
.status-3 { background: rgba(250, 204, 21, 0.2); color: #facc15; border: 1px solid rgba(250, 204, 21, 0.3); } /* No Plan - Vàng (Không gói) */
.status-4 { background: rgba(239, 68, 68, 0.2); color: #f87171; border: 1px solid rgba(239, 68, 68, 0.3); } /* No Base/Error - Đỏ */

/* SPECIFIC SIM UI */
.filter-bar { display: flex; justify-content: space-between; align-items: flex-end; gap: 20px; flex-wrap: wrap; }
.selection-bar { 
    background: var(--primary); color: white; padding: 10px 20px; border-radius: 8px; 
    display: none; align-items: center; justify-content: space-between; margin-bottom: 15px; 
    animation: slideDown 0.2s ease-out;
}
@keyframes slideDown { from { transform: translateY(-10px); opacity: 0; } to { transform: translateY(0); opacity: 1; } }

/* MODAL */
.modal-overlay { 
    position: fixed; top: 0; left: 0; width: 100%; height: 100%; 
    background: rgba(0,0,0,0.7); backdrop-filter: blur(2px);
    display: none; justify-content: center; align-items: center; z-index: 999; 
}
.modal-content { background: var(--bg-panel); border: 1px solid var(--border-color); width: 400px; border-radius: 16px; padding: 24px; box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.5); }
.modal-title { font-size: 18px; font-weight: 700; margin-bottom: 20px; color: var(--text-main); }
`;

const STATE = {
    activeTab: "configs",
    baseUrl: "http://localhost:8000",
    sims: [],
    selectedIds: new Set(),
    pagination: { page: 1, limit: 10, total: 0 }
};

const getStatusLabel = (status) => {
    switch (status) {
        case 1: return '<span class="badge status-1">Normal (Mới)</span>';
        case 2: return '<span class="badge status-2">Has Plan (Có gói)</span>';
        case 3: return '<span class="badge status-3">No Plan (Trống)</span>';
        case 4: return '<span class="badge status-4">Error (Lỗi)</span>';
        default: return `<span class="badge">Unknown (${status})</span>`;
    }
};

const apiCall = async (path, method = "GET", body = null) => {
    try {
        const options = {
            method,
            headers: {}
        };

        // FIX: Chỉ set Header và Body khi thực sự có dữ liệu gửi đi
        if (body) {
            if (body instanceof FormData) {
                // Để browser tự set Content-Type cho multipart (upload file)
                options.body = body;
            } else {
                // JSON bình thường
                options.headers["Content-Type"] = "application/json";
                options.body = JSON.stringify(body);
            }
        }

        const res = await fetch(`${STATE.baseUrl}${path}`, options);
        
        // Xử lý response 204 (No Content) - thường gặp ở DELETE/OPTIONS
        if (res.status === 204) return { ok: true };

        const data = await res.json();
        return { ok: res.ok, data };

    } catch (e) { 
        console.error("API Call Error:", e); 
        return { ok: false, error: e.message }; 
    }
};

// --- VIEWS ---

const ConfigView = () => `
    <div style="display:flex; justify-content:space-between; align-items:center; margin-bottom:20px">
        <h2 style="margin:0">System Configuration</h2>
        <button onclick="fetchQueue()">↻ Refresh Monitor</button>
    </div>

    <div id="queue-container" class="grid" style="margin-bottom:30px"></div>

    <div class="panel">
        <h3 style="margin-top:0; margin-bottom:20px; color:var(--text-muted)">Global Settings</h3>
        <div id="config-container" class="grid"></div>
    </div>
`;

const SimView = () => `
    <div class="panel">
        <div class="filter-bar">
            <div style="display:flex; gap:12px; align-items: flex-end;">
                <div class="form-group" style="margin:0">
                    <label>Filter Note</label>
                    <input id="f-note" placeholder="Search note..." style="width:200px" onkeydown="if(event.key==='Enter') fetchSims(1)">
                </div>
                <div class="form-group" style="margin:0">
                    <label>Status</label>
                    <select id="f-status" style="width:140px" onchange="fetchSims(1)">
                        <option value="">All Status</option>
                        <option value="1">Normal (1)</option>
                        <option value="2">Has Plan (2)</option>
                        <option value="3">No Plan (3)</option>
                        <option value="4">Error (4)</option>
                    </select> 
                </div>
                <button onclick="fetchSims(1)">Search</button>
            </div>
            
            <div style="display:flex; gap:8px">
                <button class="primary" onclick="triggerSyncAll()">▶ SYNC ALL DB</button>
                <button onclick="openModal()">+ Create</button>
                <button onclick="document.getElementById('file-import').click()">↑ Import</button>
                <button onclick="exportExcel()">↓ Export</button>
                <input type="file" id="file-import" hidden onchange="handleImport(this)">
            </div>
        </div>
    </div>

    <div id="selection-bar" class="selection-bar">
        <span id="selection-count" style="font-weight:600">0 items selected</span>
        <div style="display:flex; gap:10px">
            <button style="background: rgba(255,255,255,0.2); color: white; border:none" onclick="bulkSync()">⚡ Sync Selected</button>
            <button style="background: rgba(239, 68, 68, 1); color: white; border:none" onclick="bulkDelete()">🗑 Delete Selected</button>
        </div>
    </div>

    <div class="panel" style="padding:0; overflow:hidden">
        <table>
            <thead>
                <tr>
                    <th style="width:40px; text-align:center"><input type="checkbox" onchange="toggleAll(this.checked)"></th>
                    <th>ISDN</th>
                    <th>Serial</th>
                    <th>Status</th>
                    <th>Note</th>
                    </tr>
            </thead>
            <tbody id="sim-table-body"></tbody>
        </table>
        
        <div id="pagination" style="padding:16px; display:flex; justify-content:space-between; align-items:center; border-top: 1px solid var(--border-color); background: rgba(0,0,0,0.1)"></div>
    </div>

    <div id="modal-create" class="modal-overlay" onclick="if(event.target===this) closeModal()">
        <div class="modal-content">
            <div class="modal-title">Create New SIM</div>
            <div class="form-group">
                <label>ISDN</label>
                <input id="m-isdn" placeholder="Example: 0912345678">
            </div>
            <div class="form-group">
                <label>Serial</label>
                <input id="m-serial" placeholder="Example: 89840xxxxxxxx">
            </div>
            <div class="form-group">
                <label>Initial Note</label>
                <input id="m-note" placeholder="Optional note">
            </div>
            <div style="display:flex; justify-content: flex-end; gap:10px; margin-top:25px">
                <button onclick="closeModal()">Cancel</button>
                <button class="primary" onclick="saveSingleSim()">Create SIM</button>
            </div>
        </div>
    </div>
`;

// --- MAIN CONTROLLER ---

window.switchTab = (tab) => {
    STATE.activeTab = tab;
    renderLayout();
    if(tab === 'configs') { fetchConfigs(); fetchQueue(); }
    else { fetchSims(1); }
};

// =======================
// CONFIG LOGIC
// =======================

async function fetchConfigs() {
    const res = await apiCall("/api/config");
    if(res.ok) {
        document.getElementById("config-container").innerHTML = res.data.map(c => `
            <div class="panel" style="background: rgba(255,255,255,0.02); margin:0; border: 1px dashed var(--border-color)">
                <div class="form-group">
                    <label>Config Code</label>
                    <input id="cfg-code-${c.ID}" value="${c.Code}" disabled>
                </div>
                <div class="form-group">
                    <label>API URL</label>
                    <input id="cfg-url-${c.ID}" value="${c.ApiUrl}">
                </div>
                <div class="form-group">
                    <label>Token</label>
                    <input id="cfg-token-${c.ID}" value="${c.Token}" type="password">
                </div>
                <div class="form-group">
                    <label>Device ID</label>
                    <input id="cfg-device-${c.ID}" value="${c.DeviceId}">
                </div>
                <div style="display:flex; gap:12px">
                    <div class="form-group" style="flex:1">
                        <label>Import Threads</label>
                        <input type="number" id="cfg-import-${c.ID}" value="${c.ImportSimConcurrency}">
                    </div>
                    <div class="form-group" style="flex:1">
                        <label>Sync Threads</label>
                        <input type="number" id="cfg-sync-${c.ID}" value="${c.SyncSimConcurrency}">
                    </div>
                </div>
                <button class="primary" style="width:100%" onclick="saveConfig(${c.ID})">SAVE CONFIGURATION</button>
            </div>
        `).join('');
    }
}

window.saveConfig = async (id) => {
    const body = {
        Code: document.getElementById(`cfg-code-${id}`).value,
        ApiUrl: document.getElementById(`cfg-url-${id}`).value,
        Token: document.getElementById(`cfg-token-${id}`).value,
        DeviceId: document.getElementById(`cfg-device-${id}`).value,
        ImportSimConcurrency: parseInt(document.getElementById(`cfg-import-${id}`).value),
        SyncSimConcurrency: parseInt(document.getElementById(`cfg-sync-${id}`).value)
    };

    const res = await apiCall(`/api/config/${id}`, "PATCH", body);
    if(res.ok) {
        alert("Config saved & Worker pools adjusted!");
        fetchQueue();
    } else {
        alert("Failed to save config");
    }
};

window.fetchQueue = async () => {
    const res = await apiCall("/api/config/queue");
    if(res.ok) {
        const { import_queue, sync_queue } = res.data;
        document.getElementById("queue-container").innerHTML = `
            <div class="panel" style="background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%); margin:0; border-top: 4px solid var(--primary); position:relative; overflow:hidden">
                <label>Import Queue</label>
                <div style="font-size:32px; font-weight:800; margin:10px 0; color:white">${import_queue.pending}</div>
                <div style="display:flex; justify-content:space-between; font-size:12px">
                    <span style="color:var(--success)">● ${import_queue.active_workers} Active</span>
                    <span style="color:var(--text-muted)">Cap: ${import_queue.capacity}</span>
                </div>
            </div>
            <div class="panel" style="background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%); margin:0; border-top: 4px solid var(--warning); position:relative; overflow:hidden">
                <label>Sync Queue</label>
                <div style="font-size:32px; font-weight:800; margin:10px 0; color:white">${sync_queue.pending}</div>
                <div style="display:flex; justify-content:space-between; font-size:12px">
                    <span style="color:var(--success)">● ${sync_queue.active_workers} Active</span>
                    <span style="color:var(--text-muted)">Cap: ${sync_queue.capacity}</span>
                </div>
            </div>
        `;
    }
};

// =======================
// SIM LOGIC
// =======================

window.fetchSims = async (p) => {
    STATE.pagination.page = p;
    const note = document.getElementById("f-note")?.value.trim();
    const status = document.getElementById("f-status")?.value;
    
    // --- KHẮC PHỤC LỖI PARAM RỖNG ---
    const params = new URLSearchParams();
    params.append('page', p);
    params.append('limit', STATE.pagination.limit);
    
    // Chỉ append nếu có giá trị thực sự
    if (note) params.append('note', note);
    if (status) params.append('status', status);

    // Gọi API với query string chuẩn
    const res = await apiCall(`/api/sim?${params.toString()}`);
    
    if(res.ok) {
        STATE.sims = res.data.items || [];
        STATE.pagination.total = res.data.total;
        STATE.selectedIds.clear(); 
        renderSimTable();
    }
};

function renderSimTable() {
    const body = document.getElementById("sim-table-body");
    if(STATE.sims.length === 0) {
        body.innerHTML = `<tr><td colspan="6" style="text-align:center; padding:30px; color:var(--text-muted)">No Data Found</td></tr>`;
        updatePagination();
        updateSelectionBar();
        return;
    }

    body.innerHTML = STATE.sims.map(s => `
        <tr>
            <td style="text-align:center">
                <input type="checkbox" value="${s.ID}" 
                ${STATE.selectedIds.has(String(s.ID)) ? 'checked' : ''} 
                onchange="toggleOne('${s.ID}')">
            </td>
            <td style="font-weight:600; font-size:15px; color:#e2e8f0">${s.Isdn}</td>
            <td style="font-family:monospace; color:var(--text-muted)">${s.Serial}</td>
            
            <td>${getStatusLabel(s.status)}</td>
            
            <td style="font-size:13px; color:var(--text-muted); max-width:250px; overflow:hidden; text-overflow:ellipsis; white-space:nowrap" 
                title="${s.note || ''}">
                ${s.note || '-'}
            </td>
            
        </tr>
    `).join('');
    
    updatePagination();
    updateSelectionBar();
}

// --- Checkbox Logic ---
window.toggleOne = (id) => {
    STATE.selectedIds.has(id) ? STATE.selectedIds.delete(id) : STATE.selectedIds.add(id);
    updateSelectionBar();
};
window.toggleAll = (checked) => {
    if(checked) STATE.sims.forEach(s => STATE.selectedIds.add(String(s.ID)));
    else STATE.selectedIds.clear();
    renderSimTable(); // Re-render checkboxes
    updateSelectionBar();
};
function updateSelectionBar() {
    const bar = document.getElementById("selection-bar");
    if(STATE.selectedIds.size > 0) {
        bar.style.display = "flex";
        document.getElementById("selection-count").innerText = `${STATE.selectedIds.size} Items Selected`;
    } else {
        bar.style.display = "none";
    }
}

// --- Bulk Actions ---
window.bulkDelete = async () => {
    const ids = Array.from(STATE.selectedIds);
    if(!confirm(`Are you sure you want to delete ${ids.length} items? This cannot be undone.`)) return;
    
    // Gọi DELETE với Body { IDs: [...] }
    const res = await apiCall("/api/sim", "DELETE", { ids: ids });
    if(res.ok) {
        alert("Deleted successfully!");
        fetchSims(STATE.pagination.page);
    } else {
        alert("Delete failed!");
    }
};

window.bulkSync = async () => {
    const ids = Array.from(STATE.selectedIds);
    // Gọi API Sync Batch (Bác cần implement API này bên Go)
    alert(`Syncing ${ids.length} items... (Check Logs)`);
    const res = await apiCall("/api/sim/sync-sim/list", "POST", { ids: ids });
    if(res.ok) {
        STATE.selectedIds.clear();
        fetchSims(STATE.pagination.page);
    }
};

// --- Single Actions ---
window.triggerSyncAll = async () => {
    if(!confirm("Warning: This will scan the ENTIRE database. Continue?")) return;
    const res = await apiCall("/api/sim/sync-sim", "POST");
    if(res.ok) alert("Sync All Process Started in Background.");
};

window.syncOne = async (id) => {
    const res = await apiCall(`/api/sim/sync-sim/${id}`, "POST");
    if(res.ok) fetchSims(STATE.pagination.page);
};

// --- Import / Export ---
window.exportExcel = () => {
    const note = document.getElementById("f-note").value.trim();
    const status = document.getElementById("f-status").value;

    const params = new URLSearchParams();
    if (note) params.append('note', note);
    if (status) params.append('status', status);

    window.location.href = `${STATE.baseUrl}/api/sim/export?${params.toString()}`;
};

window.handleImport = async (input) => {
    if(!input.files.length) return;
    const formData = new FormData();
    formData.append("file", input.files[0]);

    const res = await apiCall("/api/sim/import", "POST", formData);
    if(res.ok) {
        alert("File Uploaded! Processing in background...");
        fetchSims(1);
    } else {
        alert("Import Failed!");
    }
    input.value = ""; // Reset input
};

// --- Modal Create ---
window.openModal = () => document.getElementById("modal-create").style.display = "flex";
window.closeModal = () => document.getElementById("modal-create").style.display = "none";

window.saveSingleSim = async () => {
    const body = {
        Isdn: document.getElementById("m-isdn").value,
        Serial: document.getElementById("m-serial").value,
        Note: document.getElementById("m-note").value,
        Status: 1
    };
    if(!body.Isdn || !body.Serial) return alert("ISDN and Serial are required!");

    const res = await apiCall("/api/sim", "POST", body); // Đã sửa lại thành /api/sim
    if(res.ok) {
        closeModal();
        // Clear inputs
        document.getElementById("m-isdn").value = "";
        document.getElementById("m-serial").value = "";
        document.getElementById("m-note").value = "";
        fetchSims(1);
    } else {
        alert(res.data.error || "Creation failed");
    }
};

function updatePagination() {
    const totalPages = Math.ceil(STATE.pagination.total / STATE.pagination.limit);
    document.getElementById("pagination").innerHTML = `
        <div style="color:var(--text-muted); font-size:13px">
            Total: <b>${STATE.pagination.total}</b> records | Page <b>${STATE.pagination.page}</b> of ${totalPages || 1}
        </div>
        <div style="display:flex; gap:8px">
            <button onclick="fetchSims(${STATE.pagination.page - 1})" ${STATE.pagination.page <= 1 ? 'disabled' : ''}>Prev</button>
            <button onclick="fetchSims(${STATE.pagination.page + 1})" ${STATE.pagination.page >= totalPages ? 'disabled' : ''}>Next</button>
        </div>
    `;
}

// --- INITIALIZATION ---
function renderLayout() {
    document.getElementById("root").innerHTML = `
        <style>${STYLES}</style>
        <div class="app-container">
            <aside class="sidebar">
                <div class="sidebar-header">RECON DASHBOARD</div>
                <div class="nav-item ${STATE.activeTab === 'configs' ? 'active' : ''}" onclick="switchTab('configs')">Configurations</div>
                <div class="nav-item ${STATE.activeTab === 'sims' ? 'active' : ''}" onclick="switchTab('sims')">SIM Management</div>
                <div style="margin-top:auto; padding:20px; font-size:11px; color:var(--text-muted); border-top:1px solid var(--border-color)">
                    System Ready<br>v2.0.0
                </div>
            </aside>
            <main class="main-content">
                <div class="view-area">${STATE.activeTab === 'configs' ? ConfigView() : SimView()}</div>
            </main>
        </div>
    `;
}

renderLayout();
switchTab('configs');