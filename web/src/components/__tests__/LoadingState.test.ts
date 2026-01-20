import { mount } from '@vue/test-utils';
import LoadingState from '../LoadingState.vue';

describe('LoadingState', () => {
  it('renders list skeletons by default', () => {
    const wrapper = mount(LoadingState, {
      global: {
        stubs: {
          'v-skeleton-loader': {
            template: '<div class="skeleton" :data-type="type"></div>',
            props: ['type']
          }
        }
      }
    });
    const skeletons = wrapper.findAll('.skeleton');
    expect(skeletons.length).toBe(3); // default count
  });

  it('renders specified count of skeletons', () => {
    const wrapper = mount(LoadingState, {
      props: { count: 5 },
      global: {
        stubs: {
          'v-skeleton-loader': {
            template: '<div class="skeleton"></div>',
            props: ['type']
          }
        }
      }
    });
    expect(wrapper.findAll('.skeleton').length).toBe(5);
  });

  it('renders card type skeletons', () => {
    const wrapper = mount(LoadingState, {
      props: { type: 'cards', count: 2 },
      global: {
        stubs: {
          'v-skeleton-loader': {
            template: '<div class="skeleton" :data-type="type"></div>',
            props: ['type']
          }
        }
      }
    });
    const skeletons = wrapper.findAll('.skeleton');
    expect(skeletons.length).toBe(2);
    expect(skeletons[0].attributes('data-type')).toBe('card');
  });

  it('renders table type skeletons', () => {
    const wrapper = mount(LoadingState, {
      props: { type: 'table', count: 2 },
      global: {
        stubs: {
          'v-skeleton-loader': {
            template: '<div class="skeleton" :data-type="type"></div>',
            props: ['type']
          }
        }
      }
    });
    // Should have heading + rows
    const skeletons = wrapper.findAll('.skeleton');
    expect(skeletons.length).toBe(3); // 1 heading + 2 rows
  });
});
