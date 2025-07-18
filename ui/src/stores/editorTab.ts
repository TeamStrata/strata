import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { watch } from 'vue';

export class QueryEditorTab {
  //values that I need to store
  id: number;
  name: string;
  currentLiteral: string;
  savedLiteral: string;

  constructor(id: number, name: string, literal: string) {
    this.id = id
    this.name = name
    this.currentLiteral = literal
    this.savedLiteral = literal
  }

  get isSaved() {
    return this.currentLiteral == this.savedLiteral
  }
}

export const useEditorTabStore = defineStore('editorTabs', () => {

  const map = ref<Record<number, QueryEditorTab>>({})
  let current_tempID = -1; //ID value for queries that aren't managed by the db, increments negatively

  //persistence of active tab
  const currentTab = ref(null);

  const allTabs = computed(() =>
    Object.values(map.value).map(({ id, name }) => ({ id, name }))
  );

  const recycleTempID = () => {
    const tempIDs = Object.keys(map.value)
      .map(Number)
      .filter(id => id < 0);
    if (tempIDs.length === 0) {
      current_tempID = -1;
    }
  };

  // Watch for changes to map and recycle tempID if needed

  watch(
    () => Object.keys(map.value),
    () => {
      recycleTempID();
    },
    { deep: true }
  );

  const createNewQueryTab = () => {
    const id = current_tempID--;
    const name = `untitled${id}`;
    map.value[id] = new QueryEditorTab(id, name, '');
    return id;
  };

  const getTabById = (id: number) => map.value[id];

  const openQuery = (id: number, name: string, literal: string) => {
    map.value[id] = new QueryEditorTab(id, name, literal);
    return id;
  };

  const closeQuery = (id: number) => {
    delete map.value[id];
  };

  function isSaved(id: number) {
    return getTabById(id).isSaved
  }

  return { allTabs, getTabById, openQuery, closeQuery, isSaved, createNewQueryTab, currentTab }
})