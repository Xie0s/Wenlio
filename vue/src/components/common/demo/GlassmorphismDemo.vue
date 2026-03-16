<script setup lang="ts">
const scrollImages = [
  'https://images.unsplash.com/photo-1519125323398-675f0ddb6308?auto=format&fit=crop&w=1200&q=80',
  'https://images.unsplash.com/photo-1477959858617-67f85cf4f1df?auto=format&fit=crop&w=1200&q=80',
  'https://images.unsplash.com/photo-1493244040629-496f6d136cc3?auto=format&fit=crop&w=1200&q=80',
  'https://images.unsplash.com/photo-1545239351-1141bd82e8a6?auto=format&fit=crop&w=1200&q=80',
  'https://images.unsplash.com/photo-1516402707257-787c50fc3898?auto=format&fit=crop&w=1200&q=80',
]

const useCases = [
  {
    image: 'https://images.unsplash.com/photo-1518773553398-650c184e0bb3?auto=format&fit=crop&w=1200&q=80',
    tag: 'Navbar',
    tagZh: '导航栏',
    title: 'Sticky Glass Bar',
    desc: '吸附顶部时实时混合流动的色彩',
  },
  {
    image: 'https://images.unsplash.com/photo-1557683316-973673baf926?auto=format&fit=crop&w=1200&q=80',
    tag: 'Card',
    tagZh: '内容卡片',
    title: 'Floating Info Panel',
    desc: '漂浮在摄影作品之上的信息层',
  },
  {
    image: 'https://images.unsplash.com/photo-1465146633011-14f8e0781093?auto=format&fit=crop&w=1200&q=80',
    tag: 'Modal',
    tagZh: '弹窗',
    title: 'Contextual Dialog',
    desc: '模态遮蔽中依然保留背景感知',
  },
  {
    image: 'https://images.unsplash.com/photo-1513151233558-d860c5398176?auto=format&fit=crop&w=1200&q=80',
    tag: 'Badge',
    tagZh: '标签',
    title: 'Rich Overlay Badge',
    desc: '将信息层轻盈地叠加于视觉之上',
  },
]

const gradientSwatches = [
  {
    bg: 'from-sky-400 via-blue-500 to-indigo-600',
    label: '蓝色渐变',
  },
  {
    bg: 'from-rose-400 via-fuchsia-500 to-violet-600',
    label: '玫瑰紫渐变',
  },
  {
    bg: 'from-amber-400 via-orange-500 to-red-500',
    label: '日落渐变',
  },
  {
    bg: 'from-emerald-400 via-teal-500 to-cyan-600',
    label: '翠绿渐变',
  },
]
</script>

<template>
  <div class="space-y-10">

    <!-- ─── 场景 01：滚动穿越 ─── -->
    <section class="demo-stage rounded-3xl p-6">
      <div class="mb-5 flex items-start gap-3">
        <span class="glass mt-0.5 shrink-0 rounded-full px-2.5 py-0.5 text-[11px] font-bold tracking-widest">01</span>
        <div>
          <h3 class="text-base font-semibold leading-none">滚动穿越</h3>
          <p class="mt-1.5 text-sm text-muted-foreground">
            向下滚动容器，sticky <code class="rounded-md bg-muted px-1 py-0.5 font-mono text-xs">.glass</code>
            工具条会实时采样底层图像，自动混合色彩与纹理。
          </p>
        </div>
      </div>

      <!-- 滚动容器本身作为边框载体，消除中间层 overflow:hidden 对 sticky 的干扰 -->
      <div class="h-[420px] overflow-y-auto rounded-2xl border border-white/40 scroll-visual dark:border-white/10">
        <!-- Sticky glass toolbar：sticky 直接作用于滚动容器，backdrop-filter 捕获下方滚过的图像 -->
        <div class="glass sticky top-0 z-20 mx-2 mt-2 flex items-center gap-2 rounded-2xl px-4 py-2.5">
          <span class="text-[11px] font-bold tracking-widest opacity-60">GLASS</span>
          <div class="mx-2 flex gap-1">
            <button type="button" class="rounded-lg bg-white/55 px-2.5 py-1 text-xs font-medium dark:bg-white/12">全部</button>
            <button type="button" class="rounded-lg px-2.5 py-1 text-xs text-muted-foreground hover:bg-white/40 dark:hover:bg-white/8">摄影</button>
            <button type="button" class="rounded-lg px-2.5 py-1 text-xs text-muted-foreground hover:bg-white/40 dark:hover:bg-white/8">视频</button>
          </div>
          <span class="ml-auto text-xs text-muted-foreground">↓ 在此区域内向下滚动</span>
        </div>

        <div class="space-y-3 p-3 pt-3">
          <article
            v-for="(img, i) in scrollImages"
            :key="img"
            class="overflow-hidden rounded-xl border border-white/30 bg-black/10"
          >
            <img :src="img" :alt="`背景 ${i + 1}`" class="h-44 w-full object-cover" />
            <p class="px-3 py-2 text-xs text-muted-foreground">
              第 {{ i + 1 }} 张 — 玻璃层实时采样此背景的色彩与亮度
            </p>
          </article>
        </div>
      </div>
    </section>

    <!-- ─── 场景 02：一种 class，四种场景 ─── -->
    <section class="demo-stage rounded-3xl p-6">
      <div class="mb-5 flex items-start gap-3">
        <span class="glass mt-0.5 shrink-0 rounded-full px-2.5 py-0.5 text-[11px] font-bold tracking-widest">02</span>
        <div>
          <h3 class="text-base font-semibold leading-none">一种 class，四种场景</h3>
          <p class="mt-1.5 text-sm text-muted-foreground">
            相同的 <code class="rounded-md bg-muted px-1 py-0.5 font-mono text-xs">.glass</code>
            class，在不同背景与使用语境下自然适应 — 无需额外变体。
          </p>
        </div>
      </div>

      <div class="grid grid-cols-2 gap-3 sm:grid-cols-4">
        <article
          v-for="item in useCases"
          :key="item.tag"
          class="group relative aspect-[3/4] overflow-hidden rounded-2xl border border-white/30"
        >
          <!-- 背景图：absolute 填满卡片，hover 放大 -->
          <img
            :src="item.image"
            :alt="item.title"
            class="absolute inset-0 h-full w-full object-cover transition-transform duration-500 group-hover:scale-105"
          />
          <!-- 暗角渐变，提升文字可读性 -->
          <div class="absolute inset-0 bg-gradient-to-t from-black/50 via-transparent to-transparent" />

          <!-- 顶部标签：glass + absolute，backdrop-filter 捕获下方图像 -->
          <div class="glass absolute left-2 top-2 rounded-full px-2.5 py-1">
            <span class="text-[10px] font-bold tracking-wider">{{ item.tag }}</span>
          </div>

          <!-- 底部信息面板：glass + absolute -->
          <div class="glass absolute inset-x-2 bottom-2 rounded-xl p-2.5">
            <p class="text-[11px] font-semibold leading-tight">{{ item.tagZh }} — {{ item.title }}</p>
            <p class="mt-1 text-[10px] leading-snug text-muted-foreground">{{ item.desc }}</p>
          </div>
        </article>
      </div>
    </section>

    <!-- ─── 场景 03：纯色渐变之上 ─── -->
    <section class="demo-stage rounded-3xl p-6">
      <div class="mb-5 flex items-start gap-3">
        <span class="glass mt-0.5 shrink-0 rounded-full px-2.5 py-0.5 text-[11px] font-bold tracking-widest">03</span>
        <div>
          <h3 class="text-base font-semibold leading-none">纯色渐变之上</h3>
          <p class="mt-1.5 text-sm text-muted-foreground">
            没有图片也能体验玻璃质感 — <code class="rounded-md bg-muted px-1 py-0.5 font-mono text-xs">.glass</code>
            同样适用于纯色或渐变背景，漫射底层色彩形成柔和光感。
          </p>
        </div>
      </div>

      <div class="grid grid-cols-2 gap-3 sm:grid-cols-4">
        <div
          v-for="swatch in gradientSwatches"
          :key="swatch.label"
          class="relative h-40 overflow-hidden rounded-2xl"
        >
          <!-- 渐变背景：绝对铺满，作为 glass 的 backdrop-filter 源 -->
          <div :class="`bg-gradient-to-br ${swatch.bg} absolute inset-0`" />
          <!-- glass 面板：absolute 定位于渐变之上，backdrop-filter 捕获下方渐变 -->
          <div class="glass absolute inset-x-3 bottom-3 rounded-xl p-3">
            <p class="text-xs font-semibold">{{ swatch.label }}</p>
            <p class="mt-0.5 text-[10px] text-muted-foreground">.glass 漫射底层色</p>
          </div>
        </div>
      </div>
    </section>

    <!-- ─── 场景 04：模态舞台 ─── -->
    <section class="demo-stage rounded-3xl p-6">
      <div class="mb-5 flex items-start gap-3">
        <span class="glass mt-0.5 shrink-0 rounded-full px-2.5 py-0.5 text-[11px] font-bold tracking-widest">04</span>
        <div>
          <h3 class="text-base font-semibold leading-none">模态舞台</h3>
          <p class="mt-1.5 text-sm text-muted-foreground">
            在高密度内容背景上覆盖 <code class="rounded-md bg-muted px-1 py-0.5 font-mono text-xs">.glass</code>
            弹窗，阅读性与背景感知相互平衡 — 同一个 class，承载所有场景。
          </p>
        </div>
      </div>

      <div class="relative overflow-hidden rounded-2xl">
        <img
          src="https://images.unsplash.com/photo-1519389950473-47ba0277781c?auto=format&fit=crop&w=1600&q=80"
          alt="工作台背景"
          class="h-[300px] w-full object-cover"
        />
        <div class="absolute inset-0 bg-black/22 dark:bg-black/42" />

        <div class="absolute inset-0 flex items-center justify-center p-4">
          <div class="glass w-full max-w-sm rounded-2xl p-5">
            <p class="text-sm font-semibold">确认发布 v2.1.0？</p>
            <p class="mt-2 text-xs leading-relaxed text-muted-foreground">
              发布将触发灰度分流与缓存刷新，请确认回滚预案与值班人员已就绪。
            </p>
            <div class="mt-4 flex items-center justify-between">
              <span class="text-[10px] text-muted-foreground opacity-70">同一个 .glass class</span>
              <div class="flex gap-2">
                <button
                  type="button"
                  class="glass rounded-xl px-3 py-1.5 text-xs font-medium"
                >
                  稍后处理
                </button>
                <button
                  type="button"
                  class="rounded-xl bg-primary px-3 py-1.5 text-xs font-medium text-primary-foreground"
                >
                  立即发布
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

  </div>
</template>

<style scoped>
.demo-stage {
  background:
    radial-gradient(120% 120% at 0% 0%, rgb(14 165 233 / 14%), transparent 50%),
    radial-gradient(120% 120% at 100% 100%, rgb(99 102 241 / 16%), transparent 48%),
    linear-gradient(135deg, rgb(255 255 255 / 72%), rgb(255 255 255 / 40%));
}

.dark .demo-stage {
  background:
    radial-gradient(120% 120% at 0% 0%, rgb(14 165 233 / 18%), transparent 52%),
    radial-gradient(120% 120% at 100% 100%, rgb(79 70 229 / 22%), transparent 48%),
    linear-gradient(135deg, rgb(24 24 27 / 84%), rgb(9 9 11 / 68%));
}

.scroll-visual {
  scrollbar-width: thin;
}
</style>
