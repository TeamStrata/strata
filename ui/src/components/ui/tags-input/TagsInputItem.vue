<script setup>
import { reactiveOmit } from "@vueuse/core";

import { TagsInputItem, useForwardProps } from "reka-ui";
import { cn } from "@/lib/utils";
import { computed } from "vue";

const props = defineProps({
  value: { type: [String, Object], required: true },
  disabled: { type: Boolean, required: false },
  asChild: { type: Boolean, required: false },
  as: { type: [String, Object, Function], required: false },
  class: { type: null, required: false },
  color: { type: String, required: false }
});

const delegatedProps = reactiveOmit(props, "class");

const forwardedProps = useForwardProps(delegatedProps);

const customBG = computed(() => {
  console.log(props.color);
  return props.color ? { backgroundColor: "#" + props.color } : {}
})
</script>

<template>
  <TagsInputItem v-bind="forwardedProps" :style="customBG" :class="cn(
    'flex h-5 items-center rounded-md bg-secondary data-[state=active]:ring-ring data-[state=active]:ring-2 data-[state=active]:ring-offset-2 ring-offset-background',
    props.class,
  )
    ">
    <slot />
  </TagsInputItem>
</template>
