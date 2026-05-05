<script setup lang="ts">
import { ref, watch } from "vue";
import { Button } from "@/components/ui/button";

const count = ref(0);

console.log("[Counter] - 123");

watch(count, (newValue, oldValue) => {
  console.log(`[Counter] Count changed: ${oldValue} -> ${newValue}`);
  console.log(`[Counter] Timestamp: ${new Date().toISOString()}`);
  console.log(`[Counter] Current count value: ${newValue}`);
});

const increment = () => {
  count.value++;
  console.log(`[Counter] Increment clicked, new value: ${count.value}`);
};

const decrement = () => {
  count.value--;
  console.log(`[Counter] Decrement clicked, new value: ${count.value}`);
};

const reset = () => {
  const previousValue = count.value;
  count.value = 0;
  console.log(`[Counter] Reset clicked, previous value: ${previousValue}, new value: ${count.value}`);
};
</script>

<template>
  <div class="flex flex-col items-center gap-4 p-8 border rounded-lg w-fit">
    <p class="text-lg font-semibold">Count: {{ count }}</p>
    <div class="flex gap-2">
      <Button variant="outline" @click="decrement">-</Button>
      <Button @click="increment">+</Button>
      <Button variant="destructive" @click="reset">Reset</Button>
    </div>
  </div>
</template>
