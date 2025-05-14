<script setup>
import { ref, nextTick } from "vue";
import Toast, { ToastTypes } from "./Toast.vue";

// If we can't connect to the database, or while the fetch is in progress, this displays.
var queries = ref([
	{
		id: 0,
		name: "Loading...",
		literal: "Loading..."
	},
]);
var newQueryName = ref("");
var newQueryLiteral = ref("");

var addModal = ref(false);
const toastRef = ref(null);

function addQuery() {
	// submit user and pass to form
	fetch(`${window.location.origin}/query/${newQueryName.value}`, {
		method: "POST",
		headers: {
			"Content-Type": "application/sql",
		},
		body: newQueryLiteral.value,
	})
	.then(async (response) => {
		if (!response.ok) {
			toastRef.value?.showToast(
				"There was an issue with creating the account",
				ToastTypes.FAIL,
			);
			throw new Error("Something went wrong");
		}

		toastRef.value?.showToast(
			"The account was successfully created",
			ToastTypes.SUCCESS,
		);

		// Push new user to the list
		queries.value.push({
			id: await response.text(), 
			name: newQueryName.value, 
			literal: newQueryLiteral.value
		});

		// Wait for VUE DOM to update, then apply syntax highlighting
		nextTick(hljs.highlightAll);
	})
	.catch((error) => {
		console.error("Error:", error);
	});
	addModal.value = false;
}

function deleteQuery(name) {
	fetch(`${window.location.origin}/query/${name}`, {
		method: "DELETE",
	})
	.then(res => {
		if (!res.ok) {
			toastRef.value?.showToast(
				"There was an issue deleting the user",
				ToastTypes.FAIL
			);
			throw new Error("Error deleting user");
		}

		// Handle success
		toastRef.value?.showToast(
			"The user was succesfully deleted",
			ToastTypes.SUCCESS
		);

		// Reload the user list
		loadQueries()
	})
	.catch(err => {
		console.error(err);
	});
}

function loadQueries() {
	fetch(`${window.location.origin}/queries`)
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
			queries.value = await response.json() || [];

			nextTick(hljs.highlightAll);
		})
		.catch((error) => {
			console.error(error);
		});
}

loadQueries();
</script>

<template>
	<!-- add user modal -->
	<div class="static" v-if="addModal">
		<div
			class="bg-black/20 w-screen h-screen absolute left-0 top-0 backdrop-blur-xs"
			@click="addModal = false"
		></div>
		<div
			class="bg-white p-8 w-96 m-auto rounded-lg shadow-md absolute left-1/2 transform -translate-x-1/2 top-1/3 -translate-y-1/2"
		>
			<h1 class="text-2xl font-semibold text-gray-800 mb-6">Add Query</h1>
			<div>
				<label for="newQuery" class="text-gray-600">Query Name</label>
				<input
					id="newQueryName"
					v-model="newQueryName"
					class="mt-1 border border-gray-300 rounded-md p-2 w-full focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
					placeholder="Enter query name"
				/>
			</div>
			<div class="mt-4">
				<label for="newPass" class="text-gray-600">Query Literal</label>
				<input
					id="newQueryLiteral"
					v-model="newQueryLiteral"
					class="mt-1 border border-gray-300 rounded-md p-2 w-full focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
					placeholder="Enter query SQL string"
				/>
			</div>
			<div class="flex flex-row justify-between mt-6">
				<button
					class="bg-neutral-200 px-4 py-2 rounded-md cursor-pointer hover:bg-neutral-300 focus:outline-none focus:ring-2 focus:ring-gray-300"
					@click="addModal = false"
				>
					Cancel
				</button>
				<button
					class="bg-blue-500 text-white py-2 px-4 rounded-md cursor-pointer hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-200"
					@click="addQuery"
				>
					Add Query
				</button>
			</div>
		</div>
	</div>

	<Toast ref="toastRef" />

	<p class="text-gray-700 mb-4">Manage your organizations queries</p>

	<!-- list header -->
	<div class="flex flex-row justify-between items-center mb-6">
		<h1 class="text-2xl font-semibold text-gray-800">Queries</h1>

		<div class="flex flex-row pb-2 items-center space-x-3">
			<p class="text-gray-600">{{ queries?.length }} queries</p>
			<input
				placeholder="Search"
				class="p-2 rounded-lg border border-gray-300 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent w-64"
				type="text"
			/>
			<button
				@click="addModal = true"
				class="bg-blue-500 text-white py-2 px-4 rounded-lg hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-200 cursor-pointer"
			>
				Add Query
			</button>
		</div>
	</div>

	<!-- list body -->
	<ul>
		<li
			class="border-t-2 border-neutral-200 py-4 flex flex-row justify-between items-center px-5 hover:bg-gray-50 transition-colors"
			v-for="(q, index) in queries"
			:key="index"
		>
			<span class="text-gray-800 font-medium">{{ q.name }}</span>
			 <!-- HLJS SQL Highlighting -->
			<pre><code class="language-sql rounded">{{ q.literal }}</code></pre>
			<div class="flex items-center">
				<p class="text-gray-600 mr-4">ID: {{ q.id }}</p>
				<router-link class="bg-green-500 text-white py-2 px-4 rounded-md cursor-pointer hover:bg-green-600 focus:outline-none focus:ring-2 focus:ring-green-200 mr-4" :to="`/query/run/${q.id}`">
					<svg xmlns="http://www.w3.org/2000/svg" width="25px" height="25px" viewBox="0 0 30 30" class="text-gray-600" >
						<path 
							fill="#fff" 
							d="M28,10H22V4a2.0025,2.0025,0,0,0-2-2H4A2.0025,2.0025,0,0,0,2,4V20a2.0025,2.0025,0,0,0,2,2h6v6a2.0025,2.0025,0,0,0,2,2H28a2.0025,2.0025,0,0,0,2-2V12A2.0025,2.0025,0,0,0,28,10ZM4,20V4h6V20Zm18,8V12h6V28Z"
						/>
					</svg>
				</router-link>
				<button
					class="bg-red-500 text-white py-2 px-4 rounded-md cursor-pointer hover:bg-red-600 focus:outline-none focus:ring-2 focus:ring-red-200"
					@click="deleteQuery(q.name)"
				>
					<svg xmlns="http://www.w3.org/2000/svg" width="25" height="25" viewBox="0 0 20 20" class="text-gray-600">
						<path
							fill="#fff"
							d="M5 6a4 4 0 1 1 8 0a4 4 0 0 1-8 0m-3 7c0-1.113.903-2 2.009-2h6.248A5.48 5.48 0 0 0 9 14.5c0 1.303.453 2.5 1.21 3.443Q9.617 18 9 18c-1.855 0-3.583-.386-4.865-1.203C2.833 15.967 2 14.69 2 13m17 1.5a4.5 4.5 0 1 1-9 0a4.5 4.5 0 0 1 9 0m-2.646-1.146a.5.5 0 0 0-.708-.708L14.5 13.793l-1.146-1.147a.5.5 0 0 0-.708.708l1.147 1.146l-1.147 1.146a.5.5 0 0 0 .708.708l1.146-1.147l1.146 1.147a.5.5 0 0 0 .708-.708L15.207 14.5z"
						/>
					</svg>
				</button>
			</div>
		</li>
	</ul>
</template>

<style scoped>
input {
	background-color: rgb(230, 230, 230);
	margin: 1px;
}
</style>
