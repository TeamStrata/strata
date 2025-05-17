<script setup>
import router from "@/router";
import { ref, nextTick } from "vue";
import Toast, { ToastTypes } from "./Toast.vue";

var customQuery = ref({
	id: 0,
	name: "Loading...",
	literal: "Loading...",
});

// If we can't connect to the database, or while the fetch is in progress, this displays.
var table = ref([
	{
		loading: "Loading...",
	},
]);

var newQueryName = ref("");
var newQueryLiteral = ref("");

var addModal = ref(false);
const toastRef = ref(null);
// Get current query ID
const id = window.location.pathname.split("/").pop();

function updateTable() {
	// Execute the query and get the result
	fetch(`http://localhost:8080/query/${id}/execute`)
		.then(async (response) => {
			// Handle error
			if (!response.ok) {
				toasRef.value?.showToast(
					"There was an error when loading users",
					ToastTypes.FAIL,
				);
				throw new Error("Error loading users");
			}

			// If 'null' was provided, then set to an empty array
			table.value = (await response.json()) || [];
		})
		.catch((error) => {
			console.error(error);
		});
}

function getQuery() {
	// Execute the query and get the result
	fetch(`http://localhost:8080/query/${id}`)
		.then(async (response) => {
			// Handle error
			if (!response.ok) {
				toasRef.value?.showToast(
					"There was an error when loading users",
					ToastTypes.FAIL,
				);
				throw new Error("Error loading users");
			}

			// If 'null' was provided, then set to something known
			customQuery.value = (await response.json()) || {
				id: 0,
				name: "Nil",
				literal: "Nil",
			};
		})
		.catch((error) => {
			console.error(error);
		});
}

getQuery();
updateTable();
</script>

<template>
	<Toast ref="toastRef" />

	<p class="text-gray-700 mb-4">Custom Query Rows</p>

	<!-- list header -->
	<div class="flex flex-row justify-between items-center mb-6">
		<h1 class="text-2xl font-semibold text-gray-800">
			'{{ customQuery.name }}' output:
		</h1>

		<div class="flex flex-row pb-2 items-center space-x-3">
			<p class="text-gray-600">{{ table?.length }} rows</p>
			<input
				placeholder="Search"
				class="p-2 rounded-lg border border-gray-300 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent w-64"
				type="text"
			/>
			<button
				@click="router.push('/query/list')"
				class="bg-red-500 text-white py-2 px-4 rounded-lg hover:bg-red-600 focus:outline-none focus:ring-2 focus:ring-red-200 cursor-pointer"
			>
				Back
			</button>
		</div>
	</div>

	<!-- list body -->
	<table class="w-full table-auto">
		<tr class="border-b-2 border-gray-800">
			<th
				class="border border-neutral-300 py-3 px-5"
				v-for="(column, colIndex) in table[0]"
				:key="colIndex"
			>
				<pre>{{ colIndex }}</pre>
			</th>
		</tr>
		<tr class="" v-for="(row, rowIndex) in table" :key="index">
			<td
				class="border border-neutral-300 py-3 px-5"
				v-for="(column, colIndex) in row"
				:key="colIndex"
			>
				<pre>{{ column }}</pre>
			</td>
		</tr>
	</table>
</template>

<style scoped>
input {
	background-color: rgb(230, 230, 230);
	margin: 1px;
}
</style>
