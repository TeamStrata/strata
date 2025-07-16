<script setup>
import { ref, nextTick, computed } from "vue";
import Toast, { ToastTypes } from "./Toast.vue";
import { apiFetch } from "@/api/request";
import {
	ResizableHandle,
	ResizablePanel,
	ResizablePanelGroup,
} from '@/components/ui/resizable'
import DropdownMenu from "./ui/dropdown-menu/DropdownMenu.vue";
import DropdownMenuTrigger from "./ui/dropdown-menu/DropdownMenuTrigger.vue";
import DropdownMenuContent from "./ui/dropdown-menu/DropdownMenuContent.vue";
import DropdownMenuItem from "./ui/dropdown-menu/DropdownMenuItem.vue";
import { Codemirror } from "vue-codemirror";
import { sql } from "@codemirror/lang-sql";
import Separator from "./ui/separator/Separator.vue";
import { Play, Save } from "lucide-vue-next";
import Dialog from "./ui/dialog/Dialog.vue";
import DialogHeader from "./ui/dialog/DialogHeader.vue";
import DialogContent from "./ui/dialog/DialogContent.vue";
import DialogFooter from "./ui/dialog/DialogFooter.vue";
import DialogClose from "./ui/dialog/DialogClose.vue";
import DialogDescription from "./ui/dialog/DialogDescription.vue";
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import DialogTitle from "./ui/dialog/DialogTitle.vue";
// If we can't connect to the database, or while the fetch is in progress, this displays.
var queries = ref([
	{
		id: 0,
		name: "Loading...",
		literal: "Loading...",
	},
]);
var newQueryName = ref("");

var isSaveOpen = ref(false);
const toastRef = ref(null);

function openSaveDialog() {
	isSaveOpen.value = true;

	nextTick(hljs.highlightAll);
}

function addQuery() {
	// submit user and pass to form
	apiFetch(`/query/${newQueryName.value}`, "POST", code.value, "application/sql")
		.then(async (response) => {
			if (!response.ok) {
				toastRef.value?.showToast(
					"There was an issue with creating the account",
					ToastTypes.FAIL,
				);
				throw new Error("Something went wrong");
			}

			toastRef.value?.showToast(
				"Query Saved",
				ToastTypes.SUCCESS,
			);

			// Wait for VUE DOM to update, then apply syntax highlighting
			// nextTick(hljs.highlightAll);
		})
		.catch((error) => {
			console.error("Error:", error);
		});
}

function deleteQuery(name) {
	apiFetch(`/query/${name}`, 'DELETE')
		.then((res) => {
			if (!res.ok) {
				toastRef.value?.showToast(
					"There was an issue deleting the user",
					ToastTypes.FAIL,
				);
				throw new Error("Error deleting user");
			}

			// Handle success
			toastRef.value?.showToast(
				"The user was succesfully deleted",
				ToastTypes.SUCCESS,
			);

			// Reload the user list
			loadQueries();
		})
		.catch((err) => {
			console.error(err);
		});
}

function loadQueries() {
	apiFetch('/queries', 'GET')
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
			queries.value = (await response.json()) || [];

			nextTick(hljs.highlightAll);
		})
		.catch((error) => {
			console.error(error);
		});
}

loadQueries();

function executeQuery() {
	//clear old response data
	queryResult.value = null;
	queryError.value = null;

	apiFetch('/query/executeLiteral', 'POST', code.value, 'application/sql')
		.then(async (response) => {
			if (!response.ok) {
				queryError.value = await response.text();
				throw new Error("Error executing query");
			}
			queryResult.value = await response.json();
		})
		.catch((error) => {
			console.error(error);
		});
}

function openQuery(query) {
	code.value = query.literal
}
//editor stuff
const extensions = [sql()];
const code = ref('SELECT * FROM users;');
const queryResult = ref(null);
const resultRows = computed(() => queryResult.value != null ? queryResult.value.length : null);
const queryError = ref(null);

//explorer stuff
const querySearch = ref("");
const filteredQueries = computed(() => {
	if (!querySearch.value) return queries.value;
	return queries.value.filter(q =>
		q.name.toLowerCase().includes(querySearch.value.toLowerCase())
	);
});

//chat stuff
const chatMessage = ref('');
const chatBox = ref(null);

function scaleChatBox() {
	chatBox.value.style.height = 'auto'
	chatBox.value.style.height = chatBox.value.scrollHeight + 'px'
}
</script>

<template>

	<!-- save query dialog -->
	<Dialog :open="isSaveOpen" @update:open="isSaveOpen = $event">
		<DialogContent class="sm:max-w-[425px]">
			<DialogHeader>
				<DialogTitle>Save Query</DialogTitle>
				<DialogDescription>
					Click 'Confirm' to save changes.
				</DialogDescription>
			</DialogHeader>
			<div class="grid gap-4 py-4">
				<div class="grid grid-cols-4 items-center gap-4">
					<Label for="edit-name-input" class="text-right">Name</Label>
					<Input id="edit-name-input" type="text" v-model="newQueryName" class="col-span-3" />
				</div>
				<pre><code class="language-sql rounded">{{ code }}</code></pre>
			</div>
			<DialogFooter>
				<DialogClose as-child>
					<Button type="button" variant="outline">
						Cancel
					</Button>
				</DialogClose>
				<Button type="submit" @click="addQuery">Confirm</Button>
			</DialogFooter>
		</DialogContent>
	</Dialog>

	<Toast ref="toastRef" />
	<ResizablePanelGroup id="demo-group-1" direction="horizontal" class="h-full">
		<ResizablePanel id="demo-panel-1" :default-size="20" :max-size="35" :min-size="15" collapsible
			:collapsed-size="0">
			<div class="h-full flex flex-col">
				<div class="p-2 flex justify-between bg-accent">
					<p>AI Chat</p>
				</div>
				<Separator></Separator>
				<div class="flex-2">
					CHAT GOES HERE
				</div>
				<Separator></Separator>
				<div class="py-3 px-4 flex items-center">
					<textarea placeholder="Write a message..." v-model="chatMessage" class="w-full resize-none break-words focus:outline-none"
						rows="1" ref="chatBox" @input="scaleChatBox"></textarea>
					<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none"
						stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
						class="lucide lucide-send-horizontal-icon lucide-send-horizontal" :class="chatMessage == '' ? 'text-neutral-400' : 'text-black cursor-pointer'">
						<path
							d="M3.714 3.048a.498.498 0 0 0-.683.627l2.843 7.627a2 2 0 0 1 0 1.396l-2.842 7.627a.498.498 0 0 0 .682.627l18-8.5a.5.5 0 0 0 0-.904z" />
						<path d="M6 12h16" />
					</svg>
				</div>
			</div>
		</ResizablePanel>
		<ResizableHandle id="demo-handle-2" />
		<ResizablePanel id="demo-panel-2" :default-size="50">
			<ResizablePanelGroup id="demo-group-2" direction="vertical">
				<ResizablePanel id="demo-panel-3" :default-size="50" :min-size="20">
					<div class="h-full items-center justify-center">
						<div class="p-2 flex justify-between bg-accent">
							<p>Name</p>
							<div class="flex gap-2 px-2">
								<Save @click="openSaveDialog" class="hover:cursor-pointer"></Save>
								<Play @click="executeQuery" class="hover:cursor-pointer"></Play>
							</div>
						</div>
						<Separator></Separator>
						<Codemirror v-model="code" :extensions="extensions"></Codemirror>
					</div>
				</ResizablePanel>
				<ResizableHandle id="demo-handle-2" />
				<ResizablePanel id="demo-panel-4" :default-size="50" :min-size="20">
					<div class="flex flex-col justify-center" :class="queryResult != null ? '' : 'h-full items-center'">
						<div v-if="queryResult != null">
							<div class="py-2 px-3">
								<p class="text-muted-foreground">{{ resultRows }} rows in result</p>
							</div>
							<!-- list body -->
							<table class="w-full">
								<thead>
									<tr class="bg-accent">
										<th class="border border-neutral-300 py-2 px-3 last:border-r-0 first:border-l-0 text-left"
											v-for="(column, colIndex) in queryResult[0]" :key="colIndex">
											<pre>{{ colIndex }}</pre>
										</th>
									</tr>
								</thead>
								<tbody>
									<tr v-for="(row, rowIndex) in queryResult" :key="index">
										<td class="border border-neutral-300 py-2 px-3 last:border-r-0 first:border-l-0"
											v-for="(column, colIndex) in row" :key="colIndex">
											<pre>{{ column }}</pre>
										</td>
									</tr>
								</tbody>
							</table>
						</div>
						<p v-else-if="queryError != null" class="text-destructive text-lg">{{ queryError }}</p>
						<p v-else class="text-neutral-500 text-lg my-1/4">Execute a query to see results.</p>
					</div>
				</ResizablePanel>
			</ResizablePanelGroup>
		</ResizablePanel>
		<ResizableHandle id="demo-handle-5" />
		<ResizablePanel id="demo-panel-5" :default-size="20" :max-size="35" :min-size="15" collapsible
			:collapsed-size="0">
			<div>
				<div class="p-2 flex justify-between bg-accent">
					<p>Query Explorer</p>
				</div>
				<Separator></Separator>
				<div>
					<div class="flex items-center px-3 mt-3 justify-between w-full">
						<div class="flex items-center w-full">
							<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"
								fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"
								stroke-linejoin="round"
								class="lucide lucide-search-icon lucide-search text-muted-foreground">
								<path d="m21 21-4.34-4.34" />
								<circle cx="11" cy="11" r="8" />
							</svg>
							<input v-model="querySearch" placeholder="Search for queries"
								class="px-2 focus:outline-none w-full" spellcheck="false">
						</div>
						<svg v-show="querySearch != ''" @click="querySearch = ''" xmlns="http://www.w3.org/2000/svg"
							width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor"
							stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
							class="lucide lucide-x-icon lucide-x w-4 cursor-pointer">
							<path d="M18 6 6 18" />
							<path d="m6 6 12 12" />
						</svg>
					</div>
					<ul class="w-full mt-3">
						<li v-for="(q, index) in filteredQueries" :key="index"
							class="flex items-center w-full justify-between group hover:bg-accent py-1 px-4 cursor-pointer"
							@click="openQuery(q)">
							<h1>{{ q.name }}</h1>
							<div class="opacity-0 group-hover:opacity-100">
								<DropdownMenu>
									<DropdownMenuTrigger class="cursor-pointer align-middle" @click.stop>
										<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24"
											viewBox="0 0 24 24">
											<g fill="none" stroke="currentColor" stroke-linecap="round"
												stroke-linejoin="round" stroke-width="2">
												<circle cx="12" cy="12" r="1" />
												<circle cx="19" cy="12" r="1" />
												<circle cx="5" cy="12" r="1" />
											</g>
										</svg>
									</DropdownMenuTrigger>
									<DropdownMenuContent>
										<DropdownMenuItem class="text-destructive focus:text-destructive"
											:onClick="() => { deleteQuery(q.name) }">
											<span>Delete</span>
										</DropdownMenuItem>
									</DropdownMenuContent>
								</DropdownMenu>
							</div>

						</li>
					</ul>
				</div>
			</div>
		</ResizablePanel>

	</ResizablePanelGroup>

	<!-- list body -->
	<ul class="hidden">
		<li class="border-t-2 border-neutral-200 py-4 flex flex-row justify-between items-center px-5 hover:bg-gray-50 transition-colors"
			v-for="(q, index) in queries" :key="index">
			<span class="text-gray-800 font-medium">{{ q.name }}</span>
			<!-- HLJS SQL Highlighting -->
			<pre><code class="language-sql rounded">{{ q.literal }}</code></pre>
			<div class="flex items-center">
				<p class="text-gray-600 mr-4">ID: {{ q.id }}</p>
				<router-link
					class="bg-green-500 text-white py-2 px-4 rounded-md cursor-pointer hover:bg-green-600 focus:outline-none focus:ring-2 focus:ring-green-200 mr-4"
					:to="`/query/run/${q.id}`">
					<svg xmlns="http://www.w3.org/2000/svg" width="25px" height="25px" viewBox="0 0 30 30"
						class="text-gray-600">
						<path fill="#fff"
							d="M28,10H22V4a2.0025,2.0025,0,0,0-2-2H4A2.0025,2.0025,0,0,0,2,4V20a2.0025,2.0025,0,0,0,2,2h6v6a2.0025,2.0025,0,0,0,2,2H28a2.0025,2.0025,0,0,0,2-2V12A2.0025,2.0025,0,0,0,28,10ZM4,20V4h6V20Zm18,8V12h6V28Z" />
					</svg>
				</router-link>
				<button
					class="bg-red-500 text-white py-2 px-4 rounded-md cursor-pointer hover:bg-red-600 focus:outline-none focus:ring-2 focus:ring-red-200"
					@click="deleteQuery(q.name)">
					<svg xmlns="http://www.w3.org/2000/svg" width="25" height="25" viewBox="0 0 20 20"
						class="text-gray-600">
						<path fill="#fff"
							d="M5 6a4 4 0 1 1 8 0a4 4 0 0 1-8 0m-3 7c0-1.113.903-2 2.009-2h6.248A5.48 5.48 0 0 0 9 14.5c0 1.303.453 2.5 1.21 3.443Q9.617 18 9 18c-1.855 0-3.583-.386-4.865-1.203C2.833 15.967 2 14.69 2 13m17 1.5a4.5 4.5 0 1 1-9 0a4.5 4.5 0 0 1 9 0m-2.646-1.146a.5.5 0 0 0-.708-.708L14.5 13.793l-1.146-1.147a.5.5 0 0 0-.708.708l1.147 1.146l-1.147 1.146a.5.5 0 0 0 .708.708l1.146-1.147l1.146 1.147a.5.5 0 0 0 .708-.708L15.207 14.5z" />
					</svg>
				</button>
			</div>
		</li>
	</ul>
</template>

<style scoped>
.root-wrapper {
	display: flex;
	flex-direction: row;

	.cm-editor {
		width: 0;
		flex-grow: 1;
	}
}
</style>