<script setup>
import { computed, onMounted, ref, watch } from 'vue';
import { useRoute } from 'vue-router';

import { apiFetch } from '@/api/request';
import { useCounterStore } from '@/stores/pageInfo';

import Toast, { ToastTypes } from './Toast.vue';

// Component/state data
const counterStore = useCounterStore();
const toastRef = ref(null);

// Dynamic dashboard id
const route = useRoute();
const dashboardId = computed(() => route.params.id);
watch(dashboardId, async (newId) => {
	if (newId) {
		loadDashboard(newId);
	}
});

// Dashboard data
const dashboard = ref({});

// Fetch dashboard
const loadDashboard = async (id) => {
	const url = '/dashboard/' + id;
	apiFetch(url)
	.then(async (response) => {
		if (!response.ok) {
			toastRef.value?.showToast(
				'There was an issue fetching dashboard data',
				ToastTypes.FAIL,
			);
			throw new Error('unable to fetch dashboard data');
		} else {
			dashboard.value = await response.json();
			if (dashboard.value.title) {
				counterStore.setPageInfo(dashboard.value.title, '', false);
			} else {
				toastRef.value?.showToast(
					'There was an issue fetching dashboard data',
					ToastTypes.FAIL,
				);
				throw new Error('dashboard title is missing');
			}
		}
	})
	.catch((error) => {
		console.error(error);
	});
}

onMounted(() => {
	loadDashboard(dashboardId.value);
});
</script>

<template>
	<Toast ref="toastRef" />
</template>