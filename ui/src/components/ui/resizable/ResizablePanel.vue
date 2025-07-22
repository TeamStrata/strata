<script setup>
import { SplitterPanel, useForwardPropsEmits } from "reka-ui";
import { ref, onMounted } from 'vue';


const props = defineProps({
  collapsedSize: { type: Number, required: false },
  collapsible: { type: Boolean, required: false },
  defaultSize: { type: Number, required: false },
  id: { type: String, required: false },
  maxSize: { type: Number, required: false },
  minSize: { type: Number, required: false },
  order: { type: Number, required: false },
  asChild: { type: Boolean, required: false },
  as: { type: [String, Object, Function], required: false },
});
const emits = defineEmits(["collapse", "expand", "resize"]);


const panelRef = ref(null);

defineExpose({
  collapse: () => panelRef.value?.collapse(),
  expand: () => panelRef.value?.expand(),
  isCollapsed: () => panelRef.value?.isCollapsed,
});

const forwarded = useForwardPropsEmits(props, emits);
</script>

<template>
  <SplitterPanel data-slot="resizable-panel" v-bind="forwarded" ref="panelRef">
    <slot />
  </SplitterPanel>
</template>
