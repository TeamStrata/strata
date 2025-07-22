<script setup>
import { ref, nextTick, computed, watchEffect } from "vue";
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
import { Play, Save, SaveAll, MessageSquare, Rows4, Rows3 } from "lucide-vue-next";
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
import { watch } from "vue";
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

function quickSave() {
	if (tabManager.currentTab < 0) {
		//open save dialog
		openSaveDialog();
	} else {
		//jsonify the query literal into a field called "Literal"
		const body = JSON.stringify({ Literal: code.value });
		//send a PATCH request to the server with the query name and literal
		apiFetch(`/query/${tabManager.currentTab}`, "PATCH", body, "application/json")
			.then(async (response) => {
				if (!response.ok) {
					toastRef.value?.showToast(
						"There was an issue with saving the query",
						ToastTypes.FAIL,
					);
					throw new Error("Something went wrong");
				}

				toastRef.value?.showToast(
					"Query Saved",
					ToastTypes.SUCCESS,
				);

				tabManager.getTabById(tabManager.currentTab).savedLiteral = code.value;
			})
			.catch((error) => {
				console.error("Error:", error);
			});
	}
}

function openSaveDialog() {
	isSaveOpen.value = true;

	nextTick(hljs.highlightAll);
}

function addQuery() {
	// submit user and pass to form
	apiFetch(`/query/${newQueryName.value}`, "POST", tabManager.getTabById(tabManager.currentTab).currentLiteral, "application/sql")
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

			loadQueries();
			const oldId = tabManager.currentTab;
			const newId = parseInt(await response.text());
			tabManager.reassignTabId(oldId, newId);
			focusTab(newId);
			tabManager.saveCurrentTab();
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
				toastRef.value?.showToast(
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
	let id = tabManager.openQuery(query.id, query.name, query.literal);
	focusTab(id);
}

//editor stuff
const extensions = [sql()];
const code = ref('');
const queryResult = ref(null);
const resultRows = computed(() => queryResult.value != null ? queryResult.value.length : null);
const queryError = ref(null);

//editor tabs
import { useEditorTabStore } from "@/stores/editorTab";
const tabManager = useEditorTabStore();
const activeTab = ref(null);
function focusTab(id) {
	pause();
	code.value = tabManager.getTabById(id).currentLiteral
	activeTab.value = id;
	tabManager.currentTab = id;
	resume();
}

function closeTab(id) {
	//check if there are any unsaved changes
	if (!tabManager.isSaved(id)) {
		if (!confirm(`You have unsaved changes, are you sure you want to close:\n\n ${tabManager.getTabById(activeTab.value).name}`)) {
			return;
		}
	}
	//close the tab
	tabManager.closeQuery(id);
}

//update tab names when query names change
watch(queries, (newQueries) => {
	//get all active tabs
	const allTabs = tabManager.allTabs.filter(tab => tab.id >= 0);
	//for each tab in all tabs, write the name of the query to the tab
	allTabs.forEach(tab => {
		const query = newQueries.find(q => q.id === tab.id);
		if (query) {
			tab.name = query.name;
		} else {
			//throw error if the query is not found
			console.error(`Query with id ${tab.id} not found in queries list`);
		}
	});
});

const { pause, resume } = watch(code, (newVal, oldVal) => {
	tabManager.getTabById(activeTab.value).currentLiteral = newVal;
});

const tabContainer = ref(null);
//tabs scroll
function handleTabScroll(e) {
	const container = e.currentTarget;
	container.scrollLeft += e.deltaY;
}

//tabs initial load
if (tabManager.allTabs.length < 1) {
	focusTab(tabManager.createNewQueryTab());
} else {
	focusTab(tabManager.currentTab);
}

//explorer stuff
const querySearch = ref("");
const isRenameOpen = ref(false);
const filteredQueries = computed(() => {
	if (!querySearch.value) return queries.value;
	return queries.value.filter(q =>
		q.name.toLowerCase().includes(querySearch.value.toLowerCase())
	);
});
const renameQueryName = ref("");
const renameQueryID = ref();
function renameQuery() {
	const body = JSON.stringify({ Name: renameQueryName.value });
	apiFetch(`/query/${renameQueryID.value}`, "PATCH", body).then(async (response) => {
		if (!response.ok) {
			toastRef.value?.showToast(
				"There was an issue renaming the query",
				ToastTypes.FAIL,
			);
			throw new Error("Error renaming query");
		}

		toastRef.value?.showToast(
			"Query Renamed",
			ToastTypes.SUCCESS,
		);

		isRenameOpen.value = false;
		loadQueries();
	}).catch((error) => {
		console.error("Error:", error);
	});
}

const explorerWindow = ref(null);

//chat stuff
const chatMessage = ref('');
const chatBox = ref(null);
const chatWindow = ref(null);

function scaleChatBox() {
	chatBox.value.style.height = 'auto'
	chatBox.value.style.height = chatBox.value.scrollHeight + 'px'
}
</script>

<template>

	<!-- save as query dialog -->
	<Dialog :open="isSaveOpen" @update:open="isSaveOpen = $event">
		<DialogContent class="sm:max-w-[425px]">
			<DialogHeader>
				<DialogTitle>Save As</DialogTitle>
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

	<!-- rename query dialog -->
	<Dialog :open="isRenameOpen" @update:open="isRenameOpen = $event">
		<DialogContent class="sm:max-w-[425px]">
			<DialogHeader>
				<DialogTitle>Rename Query</DialogTitle>
				<DialogDescription>
					Enter a new name for this query.
				</DialogDescription>
			</DialogHeader>
			<div class="grid gap-4 py-4">
				<div class="grid grid-cols-4 items-center gap-4">
					<Label for="rename-input" class="text-right">New Name</Label>
					<Input id="rename-input" type="text" v-model="renameQueryName" class="col-span-3"
						@keydown.enter="renameQuery" />
				</div>
			</div>
			<DialogFooter>
				<DialogClose as-child>
					<Button type="button" variant="outline">Cancel</Button>
				</DialogClose>
				<Button type="submit" @click="renameQuery">Confirm</Button>
			</DialogFooter>
		</DialogContent>
	</Dialog>

	<Toast ref="toastRef" />
	<ResizablePanelGroup id="demo-group-1" direction="horizontal" class="h-full">
		<ResizablePanel id="demo-panel-1" :default-size="20" :max-size="35" :min-size="15" collapsible
			:collapsed-size="0" ref="chatWindow">
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
					<textarea placeholder="Write a message..." v-model="chatMessage"
						class="w-full resize-none break-words focus:outline-none" rows="1" ref="chatBox"
						@input="scaleChatBox"></textarea>
					<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none"
						stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
						class="lucide lucide-send-horizontal-icon lucide-send-horizontal"
						:class="chatMessage == '' ? 'text-neutral-400' : 'text-black cursor-pointer'">
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
						<div class="h-10 flex justify-between bg-accent">
							<div class="flex items-center overflow-auto hideScrollbar" ref="tabContainer"
								@wheel.prevent="handleTabScroll">
								<div v-for="tab in tabManager.allTabs" :key="tab.id"
									class="px-2 border-r border-neutral-300 h-full relative flex items-center cursor-pointer border-b-2 text-nowrap group"
									:class="activeTab === tab.id ? 'bg-white border-b-primary' : 'border-b-transparent bg-accent'"
									@click="() => { focusTab(tab.id) }">
									<p :class="tab.id < 0 ? 'italic' : ''">{{ tab.name }}</p>
									<p v-show="!tabManager.isSaved(tab.id) || tab.id < 0">*</p>
									<div class="absolute pl-1 right-2 hidden group-hover:inline bg-inherit"
										@click="() => { closeTab(tab.id) }">
										<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20"
											viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
											stroke-linecap="round" stroke-linejoin="round"
											class="lucide lucide-x-icon lucide-x">
											<path d="M18 6 6 18" />
											<path d="m6 6 12 12" />
										</svg>
									</div>
								</div>
								<div class="ml-2 p-1 cursor-pointer hover:bg-neutral-200 rounded-sm"
									@click="() => { focusTab(tabManager.createNewQueryTab()) }">
									<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"
										fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"
										stroke-linejoin="round" class="lucide lucide-plus-icon lucide-plus">
										<path d="M5 12h14" />
										<path d="M12 5v14" />
									</svg>
								</div>
							</div>
							<div class="flex gap-2 px-3 items-center">
								<Separator class="mx-2" orientation="vertical"></Separator>
								<div class="p-1 cursor-pointer hover:bg-neutral-200 rounded-sm"
									@click="chatWindow?.isCollapsed() ? chatWindow?.expand() : chatWindow?.collapse()">
									<MessageSquare :size="24"></MessageSquare>
								</div>
								<div class="p-1 cursor-pointer hover:bg-neutral-200 rounded-sm"
									@click="explorerWindow?.isCollapsed() ? explorerWindow?.expand() : explorerWindow?.collapse()">
									<Rows3 :size="24" class="hover:cursor-pointer"></Rows3>
								</div>
								<Separator class="mx-2" orientation="vertical"></Separator>
								<div class="p-1 cursor-pointer hover:bg-neutral-200 rounded-sm" @click="quickSave">
									<Save :size="24" class="hover:cursor-pointer"></Save>
								</div>
								<div class="p-1 cursor-pointer hover:bg-neutral-200 rounded-sm" @click="openSaveDialog">
									<SaveAll :size="24" class="hover:cursor-pointer"></SaveAll>
								</div>
								<div class="p-1 cursor-pointer hover:bg-neutral-200 rounded-sm" @click="executeQuery">
									<Play :size="24" class="hover:cursor-pointer"></Play>
								</div>
							</div>
						</div>
						<Separator></Separator>
						<Codemirror v-if="tabManager.allTabs.length > 0" v-model="code" :extensions="extensions">
						</Codemirror>
						<div @click="() => { focusTab(tabManager.createNewQueryTab()) }" v-else
							class="flex text-neutral-400 relative ml-2 mt-2">
							<svg xmlns="http://www.w3.org/2000/svg" width="38" height="38" viewBox="0 0 24 24"
								fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"
								stroke-linejoin="round" class="lucide lucide-corner-left-up-icon lucide-corner-left-up">
								<path d="M14 9 9 4 4 9" />
								<path d="M20 20h-7a4 4 0 0 1-4-4V4" />
							</svg>
							<p class="cursor-default text-xl my-1/4 mt-4 ml-2">Create a query to get started!</p>
						</div>
					</div>
				</ResizablePanel>
				<ResizableHandle id="demo-handle-2" />
				<ResizablePanel id="demo-panel-4" :default-size="50" :min-size="20">
					<div class="flex flex-col h-full justify-center" :class="queryResult != null ? '' : 'items-center'">
						<div v-if="queryResult != null" class="flex flex-col h-full">
							<div class="py-2 px-3">
								<p class="text-muted-foreground">{{ resultRows }} rows in result</p>
							</div>

							<!-- Horizontal scroll container -->
							<div class="overflow-x-auto h-full">
								<!-- Inner container to support sticky header + vertical scroll -->
								<div class="min-w-full h-full flex flex-col">
									<!-- Use a single table with sticky header -->
									<table class="table-auto w-full border-collapse">
										<thead class="bg-accent">
											<tr>
												<th v-for="(column, colIndex) in queryResult[0]" :key="colIndex"
													class="sticky top-0 z-10 bg-accent border first:border-l-0 last:border-r-0 border-neutral-300 py-2 px-3 text-left">
													<pre class="whitespace-pre-wrap">{{ colIndex }}</pre>
												</th>
											</tr>
										</thead>
										<tbody class="overflow-y-auto">
											<tr v-for="(row, rowIndex) in queryResult" :key="rowIndex">
												<td v-for="(column, colIndex) in row" :key="colIndex"
													class="border first:border-l-0 last:border-r-0 border-neutral-300 py-2 px-3 align-top">
													<pre class="whitespace-pre-wrap">{{ column }}</pre>
												</td>
											</tr>
										</tbody>
									</table>
								</div>
							</div>
						</div>

						<p v-else-if="queryError != null" class="text-destructive text-lg">{{ queryError }}</p>
						<p v-else class="text-neutral-500 text-lg my-1/4">Execute a query to see results.</p>
					</div>
				</ResizablePanel>
			</ResizablePanelGroup>
		</ResizablePanel>
		<ResizableHandle id="demo-handle-5" />
		<ResizablePanel id="demo-panel-5" :default-size="20" :max-size="35" :min-size="15" collapsible
			:collapsed-size="0" ref="explorerWindow">
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
										<DropdownMenuItem class="cursor-pointer"
											@click="isRenameOpen = true; renameQueryID = q.id">
											Rename
										</DropdownMenuItem>
										<DropdownMenuItem class="text-destructive focus:text-destructive cursor-pointer"
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

.hideScrollbar {
	scrollbar-width: none;
	/* Firefox */
	-ms-overflow-style: none;
	/* IE 10+ */
}

.hideScrollbar::-webkit-scrollbar {
	display: none;
	/* Chrome, Safari */
}
</style>