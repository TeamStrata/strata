<template>
  <Dialog :open="addChartDialog" v-on:update:open="() => { addChartDialog = false }">
    <DialogContent>
      <DialogTitle>
        Add a chart
      </DialogTitle>
      <select v-model="selectedChartTitle" class="w-full p-2 border border-gray-300 rounded">
        <option value="">-- Select a chart --</option>
        <option v-for="chart in savedChartTitles" :key="chart.id" :value="chart">
          {{ chart.title }}
        </option>
      </select>
      <DialogFooter>
        <DialogClose>
          <Button variant="outline">Cancel</Button>
        </DialogClose>
        <Button :onClick="() => { loadChartFromDB(); addChartDialog = false; }">Add</Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>

  <!-- dashboard config -->
  <Dialog :open="showBoardConfig" v-on:update:open="() => { showBoardConfig = false }">
    <DialogContent>
      <DialogHeader>
        <DialogTitle>Edit Dashboard</DialogTitle>
      </DialogHeader>
      <div class="flex flex-col gap-1">
        <Label>Title</Label>
        <Input v-model="boardConfig.title"></Input>
      </div>
      <div class="flex flex-col gap-1">
        <Label>Description</Label>
        <Input v-model="boardConfig.desc"></Input>
      </div>
      <Label>Role Management</Label>
      <!-- input table here -->
      <table class="w-full text-left border border-gray-300 rounded mb-2">
        <thead>
          <tr class="bg-gray-100">
            <th class="p-2">Role ID</th>
            <th class="p-2">Can View</th>
            <th class="p-2">Can Edit</th>
            <th class="p-2">Can Delete</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(role, idx) in boardConfig.rolePermissions" :key="role.roleId">
            <td class="p-2">{{ role.roleName }}</td>
            <td class="p-2">
              <input type="checkbox" v-model="role.canView" />
            </td>
            <td class="p-2">
              <input type="checkbox" v-model="role.canEdit" />
            </td>
            <td class="p-2">
              <input type="checkbox" v-model="role.canDelete" />
            </td>
          </tr>
          <tr>
            <td class="p-2">
              <select v-model="selectedRoleToAdd" class="border rounded p-1 w-full">
                <option value="">-- Add Role --</option>
                <option v-for="role in availableRoles" :key="role.id" :value="role.id">
                  {{ role.name || role.id }}
                </option>
              </select>
            </td>
            <td class="p-2"></td>
            <td class="p-2"></td>
            <td class="p-2">
              <Button variant="outline" :disabled="!selectedRoleToAdd" :onClick="addRolePermission">
                Add Role
              </Button>
            </td>
          </tr>
        </tbody>
      </table>
      <DialogFooter>
        <DialogClose>
          <Button variant="outline">Cancel</Button>
        </DialogClose>
        <Button :onClick="updateBoardConfig">Save Changes</Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>

  <div class="flex flex-col">
    <!-- Edit mode header -->
    <div v-if="editMode">
      <div class="p-2 px-4 flex justify-between items-center bg-muted">
        <div class="flex items-center gap-3">
          <Construction class="text-amber-500"></Construction>
          <p class="text-lg">Edit Mode</p>
        </div>
        <div class="flex gap-2">
          <Button variant="outline" :onClick="() => { addChartDialog = true }">
            <Plus></Plus>Add Chart
          </Button>
          <Button variant="outline" :onClick="() => { editMode = false }">Exit</Button>
          <Button :onClick="() => { saveDashboardCharts(); editMode = false; }">Save Dashboard</Button>
        </div>
      </div>
      <Separator></Separator>
    </div>
    <!-- dashboard container -->
    <div class="p-3 flex gap-3">
      <ChartWidget v-for="chart in selectedCharts" :key="chart.id" :chart-data="chart" :width="chart.size_x"
        :height="chart.size_y" :id="`chart-widget-${chart.id}`" @close="removeChart(chart.id)"
        @update:size="(event) => onSizeUpdate(chart.id, event)" :edit-mode="editMode" />
    </div>

    <!-- <div>
      <button @click="saveDashboardCharts"
        class="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700 transition">Save Layout</button>
    </div> -->
    <Toast ref="toastRef" />
  </div>
</template>

<script setup>
import { ref, onMounted, reactive, onUnmounted, onBeforeUnmount, computed } from 'vue';
import { apiFetch } from '@/api/request';
import { useRoute } from 'vue-router';
import Toast, { ToastTypes } from './Toast.vue';
import ChartWidget from './ui/ChartWidget/ChartWidget.vue';
import Swal from 'sweetalert2';
import { usePageInfoStore } from '@/stores/pageInfo';
const route = useRoute();
const dashboardId = ref(route.params.id);

const charts = ref([]);
const selectedCharts = ref([]);
const selectedChartId = ref("");
const selectedChartTitle = ref("");
const toastRef = ref(null);
const savedChartTitles = ref([])

const pageInfoStore = usePageInfoStore();

//event bus management
import { useEventBus } from '@/stores/eventBus';
const bus = useEventBus();

const allRoles = ref([])
const availableRoles = computed(() => {
  // Only include roles not already present in boardConfig.rolePermissions
  const existingRoleIds = boardConfig.rolePermissions.map(rp => rp.roleId);
  return allRoles.value.filter(role => !existingRoleIds.includes(role.id));
});
const selectedRoleToAdd = ref(null)
async function getAllRoles() {
    const response = await apiFetch('/admin/roles')
    allRoles.value = await response.json()
}

function addRolePermission() {
  if (!selectedRoleToAdd.value) return;
  // Find the role object from allRoles to get the name
  const roleObj = allRoles.value.find(role => role.id === selectedRoleToAdd.value);
  boardConfig.rolePermissions.push({
    roleId: selectedRoleToAdd.value,
    roleName: roleObj ? roleObj.name : '', // Add the name here
    canView: true,
    canEdit: false,
    canDelete: false
  });
  selectedRoleToAdd.value = null;
}

getAllRoles();

// edit mode stuff
const editMode = ref(false);
const addChartDialog = ref(false);

// settings dialog and config
const showBoardConfig = ref(false);
const boardConfig = reactive({
  title: null,
  desc: null,
  rolePermissions: []
});

function onSizeUpdate(chartId, size) {
  const chart = selectedCharts.value.find(c => c.id === chartId);
  if (chart) {
    chart.size_x = size.width;
    chart.size_y = size.height;
  }
}

const fetchSavedChartTitles = async () => {
  try {
    const response = await apiFetch('/charts')
    if (!response.ok) throw new Error('Failed to fetch chart titles')
    savedChartTitles.value = await response.json()
  } catch (err) {
    console.error('Error fetching titles:', err)
  }
}

const saveDashboardCharts = async () => {
  try {
    // Fetch existing charts from backend
    const existingRes = await apiFetch(`/dashboard/${dashboardId.value}/charts`, 'GET');
    if (!existingRes.ok) throw new Error('Failed to fetch charts from dashboard');

    const existingCharts = await existingRes.json();

    // Map chart_id to size, cast to Number to ensure matching works
    const existingChartMap = new Map(
      Array.isArray(existingCharts)
        ? existingCharts
          .filter(chart => chart?.chart_id != null)
          .map(chart => [
            Number(chart.chart_id), // ✅ Ensure consistent type
            { size_x: chart.size_x, size_y: chart.size_y }
          ])
        : []
    );

    for (let index = 0; index < selectedCharts.value.length; index++) {
      const chart = selectedCharts.value[index];
      if (!chart?.id) continue;

      const chartIdNum = Number(chart.id); // ✅ Normalize ID for map lookup

      // 🟡 Get actual rendered size from the DOM
      const widgetEl = document.getElementById(`chart-widget-${chart.id}`);
      let size_x = chart.size_x || 300;
      let size_y = chart.size_y || 300;

      if (widgetEl) {
        const rect = widgetEl.getBoundingClientRect();
        size_x = Math.round(rect.width);
        size_y = Math.round(rect.height);
      }

      const existing = existingChartMap.get(chartIdNum);

      const shouldUpdate =
        !existing ||
        existing.size_x !== size_x ||
        existing.size_y !== size_y;

      // Debug logs
      console.log(`Chart ID ${chart.id}:`);
      console.log(`  Existing:`, existing);
      console.log(`  New Size: ${size_x}x${size_y}`);
      console.log(`  shouldUpdate:`, shouldUpdate);

      if (!shouldUpdate) {
        console.log(`Chart ID ${chart.id} already exists with same size. Skipping...`);
        continue;
      }

      // Payload for PATCH
      const payload = {
        chart_order: chart.chart_order || index + 1,
        id: chart.id,
        size_x,
        size_y,
        title: chart.title || "",
        type: chart.type || ""
      };

      // Save or update chart
      const response = await apiFetch(
        `/dashboard/${dashboardId.value}/chart/${chart.id}`,
        'PATCH',
        JSON.stringify(payload)
      );

      if (!response.ok) {
        toastRef.value?.showToast(
          `Failed to save chart ${chart.name || chart.id} to dashboard`,
          ToastTypes.FAIL
        );
        return;
      }
    }

    toastRef.value?.showToast('Dashboard layout saved successfully', ToastTypes.SUCCESS);
  } catch (error) {
    console.error(error);
    toastRef.value?.showToast('Error saving dashboard layout', ToastTypes.FAIL);
  }
};

const loadDashboardGraphs = async () => {
  try {
    const response = await apiFetch(`/dashboard/${dashboardId.value}/charts`);

    if (response.status === 404) {
      return;
    }

    const graphMappings = await response.json();
    console.log(graphMappings);
    if (graphMappings != null) {
      for (const mapping of graphMappings) {
        const chartResponse = await apiFetch(`/chart/${mapping.chart_id}`);
        if (chartResponse.ok) {
          const chart = await chartResponse.json();
          chart.size_x = mapping.size_x;
          chart.size_y = mapping.size_y;
          chart.order = mapping.order;
          selectedCharts.value.push(chart);
        }
      }
    } else {
      selectedCharts.value = [];
    }
  } catch (err) {
    if (err?.response?.status !== 404) {
      toastRef.value?.showToast('Error loading dashboard graphs', ToastTypes.FAIL);
    }
    console.error(err);
  }
};

const removeChart = async (chartId) => {
  const result = await Swal.fire({
    title: 'Remove Chart?',
    text: 'Are you sure you want to remove this chart from the dashboard?',
    icon: 'warning',
    showCancelButton: true,
    confirmButtonColor: '#d33',
    cancelButtonColor: '#3085d6',
    confirmButtonText: 'Yes, remove it!',
  });

  if (!result.isConfirmed) return;

  try {
    // Remove from the local state
    selectedCharts.value = selectedCharts.value.filter(chart => chart.id !== chartId);

    // Delete from backend
    const res = await apiFetch(`/dashboard/${dashboardId.value}/chart/${chartId}`, 'DELETE');
    if (!res.ok) {
      toastRef.value?.showToast(`Failed to remove chart ${chartId} from dashboard`, ToastTypes.FAIL);
    } else {
      toastRef.value?.showToast(`Chart ${chartId} removed from dashboard`, ToastTypes.SUCCESS);
    }
  } catch (err) {
    console.error(err);
    toastRef.value?.showToast('Error removing chart from dashboard', ToastTypes.FAIL);
  }
};


const loadChartFromDB = () => {
  console.log("HELLO")
  if (!selectedChartTitle.value) return;
  const existing = selectedCharts.value.find(c => c.id === selectedChartTitle.value.id);
  if (existing) return;
  selectedCharts.value.push({
    ...selectedChartTitle.value,
    size_x: selectedCharts.value.size_x || 600,
    size_y: selectedChartTitle.value.size_y || 800,
    chart_order: selectedCharts.value.length + 1
  });
  selectedChartTitle.value = "";
};

import { watch } from 'vue';
import Separator from './ui/separator/Separator.vue';
import Button from './ui/button/Button.vue';
import { Construction, Plus } from 'lucide-vue-next';
import Dialog from './ui/dialog/Dialog.vue';
import DialogContent from './ui/dialog/DialogContent.vue';
import DialogTitle from './ui/dialog/DialogTitle.vue';
import DialogFooter from './ui/dialog/DialogFooter.vue';
import DialogClose from './ui/dialog/DialogClose.vue';
import Card from './ui/card/Card.vue';
import DialogHeader from './ui/dialog/DialogHeader.vue';
import Input from './ui/input/Input.vue';
import Label from './ui/label/Label.vue';

watch(() => route.params.id,
  (newVal, old) => {
    dashboardId.value = newVal;
    loadDashboardGraphs();
    fetchSavedChartTitles();
  });

async function loadDashboard() {
  const response = await apiFetch(`/dashboard/${dashboardId.value}`);
  if (!response.ok) return;
  const data = await response.json();
  console.log(data);
  const res2 = await apiFetch(`/dashboard/${dashboardId.value}/permissions`, 'GET')
  const data2 = await res2.json();
  boardConfig.rolePermissions = data2.permissions;

  // Update store with API info
  pageInfoStore.setPageInfo(
    data.title || "",
    data.content || "",
    data.configurable ?? true
  );
}

function updateBoardConfig() {
  const payload = {
    id: parseInt(dashboardId.value, 10),
    name: boardConfig.title,
    description: boardConfig.desc,
    permissions: Array.isArray(boardConfig.rolePermissions)
      ? boardConfig.rolePermissions.map(rp => ({
        roleId: rp.roleId,
        canView: !!rp.canView,
        canEdit: !!rp.canEdit,
        canDelete: !!rp.canDelete
      }))
      : []
  };

  apiFetch(`/dashboard`, 'PATCH', JSON.stringify(payload))
    .then(res => {
      if (!res.ok) {
        toastRef.value?.showToast('Failed to update dashboard config', ToastTypes.FAIL);
        return;
      }
      toastRef.value?.showToast('Dashboard config updated', ToastTypes.SUCCESS);
      showBoardConfig.value = false;
      loadDashboard();
    })
    .catch(err => {
      console.error(err);
      toastRef.value?.showToast('Error updating dashboard config', ToastTypes.FAIL);
    });
}

watch(showBoardConfig, (val) => {
  if (val) {
    let out = pageInfoStore.pageInfo;
    boardConfig.title = out.title;
    boardConfig.desc = out.description;
  }
});

let offAlpha
let offBeta

onMounted(async () => {
  await loadDashboardGraphs();
  await fetchSavedChartTitles();

  offAlpha = bus.on('goEdit', () => { editMode.value = true })
  offBeta = bus.on('goSettings', () => { showBoardConfig.value = true })
});

onBeforeUnmount(() => {
  offAlpha && offAlpha()
  offBeta && offBeta()
})
loadDashboard();
</script>

<style scoped>
.charts-list {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  margin-bottom: 2rem;
}

.chart-card,
.create-card {
  border: 1px solid #ccc;
  border-radius: 8px;
  padding: 1rem;
  min-width: 200px;
  background: white;
  box-shadow: 0 0 5px rgba(0, 0, 0, 0.05);
  display: flex;
  flex-direction: column;
  align-items: start;
  gap: 0.5rem;
}

.create-card {
  border: 2px dashed #ccc;
  cursor: pointer;
  justify-content: center;
  align-items: center;
  text-align: center;
}

.chart-display {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
}

.tag {
  color: white;
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 0.8rem;
}
</style>
