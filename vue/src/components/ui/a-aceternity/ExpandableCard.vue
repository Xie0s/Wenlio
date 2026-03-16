<script setup lang="ts">
/**
 * ExpandableCard.vue - 可展开卡片组件（移植自 Aceternity UI）
 * 职责：点击卡片以动画展开详情，Escape 或点击外部关闭
 * 对外暴露：ExpandableCard 组件，props: cards: CardItem[]
 * 移植来源：https://ui.aceternity.com/components/expandable-card
 */
import { ref, watch, onUnmounted, useId } from 'vue'
import { useOutsideClick } from '@/composables/useOutsideClick'

export interface CardItem {
  description: string
  title: string
  src: string
  ctaText: string
  ctaLink: string
  content: string
}

const props = defineProps<{
  cards: CardItem[]
  layout?: 'list' | 'grid'
}>()

const active = ref<CardItem | null>(null)
const id = useId()
const cardRef = ref<HTMLElement | null>(null)

const onKeyDown = (e: KeyboardEvent) => {
  if (e.key === 'Escape') active.value = null
}

watch(active, (val) => {
  if (val) {
    document.body.style.overflow = 'hidden'
    window.addEventListener('keydown', onKeyDown)
  } else {
    document.body.style.overflow = 'auto'
    window.removeEventListener('keydown', onKeyDown)
  }
})

onUnmounted(() => {
  window.removeEventListener('keydown', onKeyDown)
  document.body.style.overflow = 'auto'
})

useOutsideClick(cardRef, () => { active.value = null })
</script>

<template>
  <!-- 遮罩层 -->
  <Transition name="ec-fade">
    <div
      v-if="active"
      class="fixed inset-0 bg-black/20 h-full w-full z-10"
    />
  </Transition>

  <!-- 展开状态 -->
  <Transition name="ec-scale">
    <div v-if="active" class="fixed inset-0 grid place-items-center z-[100]">
      <!-- 移动端关闭按钮 -->
      <button
        class="flex absolute top-2 right-2 lg:hidden items-center justify-center bg-white rounded-full h-6 w-6"
        @click="active = null"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="24"
          height="24"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="h-4 w-4 text-black"
        >
          <path stroke="none" d="M0 0h24v24H0z" fill="none" />
          <path d="M18 6l-12 12" />
          <path d="M6 6l12 12" />
        </svg>
      </button>

      <!-- 展开卡片主体 -->
      <div ref="cardRef" class="w-full max-w-[500px] h-full md:h-fit md:max-h-[90%] flex flex-col bg-white dark:bg-neutral-900 sm:rounded-2xl overflow-hidden">
        <img
          :src="active.src"
          :alt="active.title"
          class="w-full h-80 object-cover object-top"
        />
        <div>
          <div class="flex justify-between items-start p-4">
            <div>
              <h3 class="font-bold text-neutral-700 dark:text-neutral-200">{{ active.title }}</h3>
              <p class="text-neutral-600 dark:text-neutral-400">{{ active.description }}</p>
            </div>
            <a
              :href="active.ctaLink"
              target="_blank"
              class="px-4 py-3 text-sm rounded-full font-bold bg-green-500 text-white"
            >
              {{ active.ctaText }}
            </a>
          </div>
          <div class="pt-4 relative px-4">
            <div
              class="text-neutral-600 text-xs md:text-sm lg:text-base h-40 md:h-fit pb-10 flex flex-col items-start gap-4 overflow-auto dark:text-neutral-400 [mask:linear-gradient(to_bottom,white,white,transparent)] [scrollbar-width:none] [-ms-overflow-style:none] [-webkit-overflow-scrolling:touch]"
              v-html="active.content"
            />
          </div>
        </div>
      </div>
    </div>
  </Transition>

  <!-- 卡片列表 - list 布局 -->
  <ul v-if="props.layout !== 'grid'" class="max-w-2xl mx-auto w-full gap-4">
    <div
      v-for="card in cards"
      :key="`card-${card.title}-${id}`"
      class="p-4 flex flex-col md:flex-row justify-between items-center hover:bg-neutral-50 dark:hover:bg-neutral-800 rounded-xl cursor-pointer transition-colors duration-200"
      @click="active = card"
    >
      <div class="flex gap-4 flex-col md:flex-row">
        <img
          :src="card.src"
          :alt="card.title"
          class="h-40 w-40 md:h-14 md:w-14 rounded-lg object-cover object-top"
        />
        <div>
          <h3 class="font-medium text-neutral-800 dark:text-neutral-200 text-center md:text-left">{{ card.title }}</h3>
          <p class="text-neutral-600 dark:text-neutral-400 text-center md:text-left">{{ card.description }}</p>
        </div>
      </div>
      <button class="px-4 py-2 text-sm rounded-full font-bold bg-gray-100 hover:bg-green-500 hover:text-white text-black mt-4 md:mt-0 transition-colors duration-200">
        {{ card.ctaText }}
      </button>
    </div>
  </ul>

  <!-- 卡片列表 - grid 布局 -->
  <ul v-else class="max-w-2xl mx-auto w-full grid grid-cols-1 md:grid-cols-2 items-start gap-4">
    <div
      v-for="card in cards"
      :key="`card-${card.title}-${id}`"
      class="p-4 flex flex-col hover:bg-neutral-50 dark:hover:bg-neutral-800 rounded-xl cursor-pointer transition-colors duration-200"
      @click="active = card"
    >
      <div class="flex gap-4 flex-col w-full">
        <img
          :src="card.src"
          :alt="card.title"
          class="h-60 w-full rounded-lg object-cover object-top"
        />
        <div class="flex justify-center items-center flex-col">
          <h3 class="font-medium text-neutral-800 dark:text-neutral-200 text-center text-base">{{ card.title }}</h3>
          <p class="text-neutral-600 dark:text-neutral-400 text-center text-base">{{ card.description }}</p>
        </div>
      </div>
    </div>
  </ul>
</template>

<style scoped>
.ec-fade-enter-active,
.ec-fade-leave-active {
  transition: opacity 0.2s ease;
}
.ec-fade-enter-from,
.ec-fade-leave-to {
  opacity: 0;
}

.ec-scale-enter-active,
.ec-scale-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}
.ec-scale-enter-from,
.ec-scale-leave-to {
  opacity: 0;
  transform: scale(0.95);
}
</style>
