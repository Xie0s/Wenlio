<script setup lang="ts">
/**
 * ExpandableCardDemo.vue - 可展开卡片演示组件
 * 职责：展示 ExpandableCard 的两种布局变体（list / grid）
 * 对外暴露：无（演示专用）
 */
import { ExpandableCard } from '@/components/ui/a-aceternity'
import type { CardItem } from '@/components/ui/a-aceternity'

const listCards: CardItem[] = [
  {
    description: 'Lana Del Rey',
    title: 'Summertime Sadness',
    src: 'https://assets.aceternity.com/demos/lana-del-rey.jpeg',
    ctaText: 'Play',
    ctaLink: 'https://ui.aceternity.com/templates',
    content: `<p>Lana Del Rey, an iconic American singer-songwriter, is celebrated for her melancholic and cinematic music style. Born Elizabeth Woolridge Grant in New York City, she has captivated audiences worldwide with her haunting voice and introspective lyrics.<br /><br />Her songs often explore themes of tragic romance, glamour, and melancholia, drawing inspiration from both contemporary and vintage pop culture. With a career that has seen numerous critically acclaimed albums, Lana Del Rey has established herself as a unique and influential figure in the music industry, earning a dedicated fan base and numerous accolades.</p>`,
  },
  {
    description: 'Babbu Maan',
    title: 'Mitran Di Chhatri',
    src: 'https://assets.aceternity.com/demos/babbu-maan.jpeg',
    ctaText: 'Play',
    ctaLink: 'https://ui.aceternity.com/templates',
    content: `<p>Babu Maan, a legendary Punjabi singer, is renowned for his soulful voice and profound lyrics that resonate deeply with his audience. Born in the village of Khant Maanpur in Punjab, India, he has become a cultural icon in the Punjabi music industry.<br /><br />His songs often reflect the struggles and triumphs of everyday life, capturing the essence of Punjabi culture and traditions. With a career spanning over two decades, Babu Maan has released numerous hit albums and singles that have garnered him a massive fan following both in India and abroad.</p>`,
  },
  {
    description: 'Metallica',
    title: 'For Whom The Bell Tolls',
    src: 'https://assets.aceternity.com/demos/metallica.jpeg',
    ctaText: 'Play',
    ctaLink: 'https://ui.aceternity.com/templates',
    content: `<p>Metallica, an iconic American heavy metal band, is renowned for their powerful sound and intense performances that resonate deeply with their audience. Formed in Los Angeles, California, they have become a cultural icon in the heavy metal music industry.<br /><br />Their songs often reflect themes of aggression, social issues, and personal struggles, capturing the essence of the heavy metal genre. With a career spanning over four decades, Metallica has released numerous hit albums and singles that have garnered them a massive fan following both in the United States and abroad.</p>`,
  },
  {
    description: 'Led Zeppelin',
    title: 'Stairway To Heaven',
    src: 'https://assets.aceternity.com/demos/led-zeppelin.jpeg',
    ctaText: 'Play',
    ctaLink: 'https://ui.aceternity.com/templates',
    content: `<p>Led Zeppelin, a legendary British rock band, is renowned for their innovative sound and profound impact on the music industry. Formed in London in 1968, they have become a cultural icon in the rock music world.<br /><br />Their songs often reflect a blend of blues, hard rock, and folk music, capturing the essence of the 1970s rock era. With a career spanning over a decade, Led Zeppelin has released numerous hit albums and singles that have garnered them a massive fan following both in the United Kingdom and abroad.</p>`,
  },
  {
    description: 'Mustafa Zahid',
    title: 'Toh Phir Aao',
    src: 'https://assets.aceternity.com/demos/toh-phir-aao.jpeg',
    ctaText: 'Play',
    ctaLink: 'https://ui.aceternity.com/templates',
    content: `<p>"Aawarapan", a Bollywood movie starring Emraan Hashmi, is renowned for its intense storyline and powerful performances. Directed by Mohit Suri, the film has become a significant work in the Indian film industry.<br /><br />The movie explores themes of love, redemption, and sacrifice, capturing the essence of human emotions and relationships. With a gripping narrative and memorable music, "Aawarapan" has garnered a massive fan following both in India and abroad, solidifying Emraan Hashmi's status as a versatile actor.</p>`,
  },
]

const gridCards: CardItem[] = [
  {
    description: 'Lana Del Rey',
    title: 'Summertime Sadness',
    src: 'https://assets.aceternity.com/demos/lana-del-rey.jpeg',
    ctaText: 'Visit',
    ctaLink: 'https://ui.aceternity.com/templates',
    content: `<p>Lana Del Rey, an iconic American singer-songwriter, is celebrated for her melancholic and cinematic music style. Born Elizabeth Woolridge Grant in New York City, she has captivated audiences worldwide with her haunting voice and introspective lyrics.<br /><br />Her songs often explore themes of tragic romance, glamour, and melancholia, drawing inspiration from both contemporary and vintage pop culture. With a career that has seen numerous critically acclaimed albums, Lana Del Rey has established herself as a unique and influential figure in the music industry, earning a dedicated fan base and numerous accolades.</p>`,
  },
  {
    description: 'Babbu Maan',
    title: 'Mitran Di Chhatri',
    src: 'https://assets.aceternity.com/demos/babbu-maan.jpeg',
    ctaText: 'Visit',
    ctaLink: 'https://ui.aceternity.com/templates',
    content: `<p>Babu Maan, a legendary Punjabi singer, is renowned for his soulful voice and profound lyrics that resonate deeply with his audience. Born in the village of Khant Maanpur in Punjab, India, he has become a cultural icon in the Punjabi music industry.<br /><br />His songs often reflect the struggles and triumphs of everyday life, capturing the essence of Punjabi culture and traditions. With a career spanning over two decades, Babu Maan has released numerous hit albums and singles that have garnered him a massive fan following both in India and abroad.</p>`,
  },
  {
    description: 'Metallica',
    title: 'For Whom The Bell Tolls',
    src: 'https://assets.aceternity.com/demos/metallica.jpeg',
    ctaText: 'Visit',
    ctaLink: 'https://ui.aceternity.com/templates',
    content: `<p>Metallica, an iconic American heavy metal band, is renowned for their powerful sound and intense performances that resonate deeply with their audience. Formed in Los Angeles, California, they have become a cultural icon in the heavy metal music industry.<br /><br />Their songs often reflect themes of aggression, social issues, and personal struggles, capturing the essence of the heavy metal genre. With a career spanning over four decades, Metallica has released numerous hit albums and singles that have garnered them a massive fan following both in the United States and abroad.</p>`,
  },
  {
    description: 'Lord Himesh',
    title: 'Aap Ka Suroor',
    src: 'https://assets.aceternity.com/demos/aap-ka-suroor.jpeg',
    ctaText: 'Visit',
    ctaLink: 'https://ui.aceternity.com/templates',
    content: `<p>Himesh Reshammiya, a renowned Indian music composer, singer, and actor, is celebrated for his distinctive voice and innovative compositions. Born in Mumbai, India, he has become a prominent figure in the Bollywood music industry.<br /><br />His songs often feature a blend of contemporary and traditional Indian music, capturing the essence of modern Bollywood soundtracks. With a career spanning over two decades, Himesh Reshammiya has released numerous hit albums and singles that have garnered him a massive fan following both in India and abroad.</p>`,
  },
]
</script>

<template>
  <div class="space-y-12">
    <!-- 列表布局 -->
    <section class="space-y-4">
      <div class="rounded-2xl border bg-card/40 p-4">
        <h3 class="text-sm font-semibold">列表布局（List）</h3>
        <p class="text-sm text-muted-foreground mt-1">横排卡片，点击展开详情，支持 Escape / 点击外部关闭。</p>
      </div>
      <ExpandableCard :cards="listCards" />
    </section>

    <div class="border-t border-dashed" />

    <!-- 网格布局 -->
    <section class="space-y-4">
      <div class="rounded-2xl border bg-card/40 p-4">
        <h3 class="text-sm font-semibold">网格布局（Grid）</h3>
        <p class="text-sm text-muted-foreground mt-1">双列网格，大图展示，点击展开详情。</p>
      </div>
      <ExpandableCard :cards="gridCards" layout="grid" />
    </section>
  </div>
</template>
