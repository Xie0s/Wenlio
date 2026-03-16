import type { InjectionKey, Ref } from "vue"

export const TABS_SLIDE_DIRECTION_KEY = Symbol('tabs-slide-direction') as InjectionKey<Ref<'left' | 'right'>>

export const TABS_REGISTERED_VALUES_KEY = Symbol('tabs-registered-values') as InjectionKey<{
  register: (val: string) => void
  getDirection: () => 'left' | 'right'
}>
