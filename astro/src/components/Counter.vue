<script setup lang="ts">
import { ref } from "vue";
import { Button } from "@/components/ui/button";

const count = ref(0);
</script>

<template>
  <div class="flex flex-col items-center gap-4 p-8 border rounded-lg w-fit">
    <p class="text-lg font-semibold">Count: {{ count }}</p>
    <div class="flex gap-2">
      <Button variant="outline" @click="count--">-</Button>
      <Button @click="count++">+</Button>
      <Button variant="destructive" @click="count = 0">Reset</Button>
    </div>
  </div>
</template>
