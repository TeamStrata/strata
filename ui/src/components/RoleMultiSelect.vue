<script setup lang="ts">
import { useFilter } from 'reka-ui'
import { computed, ref } from 'vue'
import { Combobox, ComboboxAnchor, ComboboxEmpty, ComboboxGroup, ComboboxInput, ComboboxItem, ComboboxList } from '@/components/ui/combobox'
import { TagsInput, TagsInputInput, TagsInputItem, TagsInputItemDelete, TagsInputItemText } from '@/components/ui/tags-input'

export interface Role {
  id: number;
  name: string;
  color: string;
  usercount: number;
}

const props = defineProps<{
  allVals: Role[]
  activeRoles: string[]
}>();

const modelValue = ref<Role[]>([])
const open = ref(false)
const searchTerm = ref('')

const { contains } = useFilter({ sensitivity: 'base' })

//compute list of detailed roles that the user has
const activeRolesFull = computed(() => {
  const roleMap = new Map(props.allVals.map(role => [role.name, role]));
  return props.activeRoles.map(id => roleMap.get(id)).filter(Boolean) as Role[];
})

//compute all roles that the user does not have
const availableRoles = computed(() => {
  const selectedIds = new Set(activeRolesFull.value.map(role => role.id));
  const options = props.allVals.filter(role => !selectedIds.has(role.id));
  return searchTerm.value
    ? options.filter(option => contains(option.name, searchTerm.value))
    : options;
})

const emits = defineEmits(['roleAdd', 'roleRemove']);
</script>

<template>
  <Combobox v-model="modelValue" v-model:open="open" :ignore-filter="true">
    <ComboboxAnchor as-child>
      <TagsInput v-model="modelValue" class="px-2 gap-2 w-80">
        <div class="flex gap-2 flex-wrap items-center">
          <TagsInputItem v-for="item in activeRolesFull" :key="item.id" :value="item.name" :color="item.color">
            <TagsInputItemText />
            <TagsInputItemDelete/>
          </TagsInputItem>
        </div>

        <ComboboxInput v-model="searchTerm" as-child>
          <TagsInputInput placeholder="Add Role..."
            class="min-w-[200px] w-full p-0 border-none focus-visible:ring-0 h-auto" @keydown.enter.prevent />
        </ComboboxInput>
      </TagsInput>

      <ComboboxList class="w-[--reka-popper-anchor-width]">
        <ComboboxEmpty />
        <ComboboxGroup>
          <ComboboxItem v-for="val in availableRoles" :key="val.id" :value="val.name" :color="val.color" @select.prevent="() => {
            $emit('roleAdd', val.id);
          }">
            {{ val.name }}
          </ComboboxItem>
        </ComboboxGroup>
      </ComboboxList>
    </ComboboxAnchor>
  </Combobox>
</template>