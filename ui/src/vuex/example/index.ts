import { InjectionKey } from 'vue';
import { createStore, useStore as aliasUseStore, Store } from 'vuex';

interface exmaple {
  name: string;
  age: number;
}
const key: InjectionKey<Store<exmaple>> = Symbol();
const registStore = () => {
  return aliasUseStore(key);
};
const createThatStore = createStore<exmaple>({
  state: {
    name: 'example_name',
    age: 16
  },
  getters: {
    birthday(state): number {
      return new Date().getFullYear() - state.age;
    }
  },
  mutations: {
    addAge(state, p: { num: number }) {
      state.age += p.num;
    }
  },
  actions: {}
});
export { key, registStore, exmaple, createThatStore };
