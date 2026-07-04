import { mount } from '@vue/test-utils';
import ItemWarning from '../ItemWarning.vue';

/* jscpd:ignore-start */
const iconStubs = {
  'v-tooltip': {
    template: '<div class="tooltip"><slot name="activator" :props="{}" /><slot /></div>',
    props: ['location', 'contentClass'],
  },
  'v-btn': {
    template: '<button :data-cy="$attrs[\'data-cy\']"><slot /></button>',
    props: ['icon', 'variant', 'size', 'color'],
  },
  'v-icon': { template: '<i class="icon"><slot /></i>', props: ['size'] },
};
/* jscpd:ignore-end */

describe('ItemWarning', () => {
  it('renders the icon with the warning message on hover (icon mode)', () => {
    const wrapper = mount(ItemWarning, {
      props: { warning: 'Apple Only, Thunderbolt Display', dataCy: 'folded-warning-icon' },
      global: { stubs: iconStubs },
    });
    // The icon button carries the passed data-cy.
    expect(wrapper.find('[data-cy="folded-warning-icon"]').exists()).toBe(true);
    // The warning text is rendered (in the tooltip content slot).
    expect(wrapper.text()).toContain('Apple Only, Thunderbolt Display');
  });

  it('renders an inline styled message (inline mode)', () => {
    const wrapper = mount(ItemWarning, {
      props: { warning: 'No monitor', mode: 'inline', dataCy: 'item-warning' },
      global: { stubs: { 'v-icon': { template: '<i class="icon"><slot /></i>' } } },
    });
    const block = wrapper.find('[data-cy="item-warning"]');
    expect(block.exists()).toBe(true);
    expect(block.classes()).toContain('item-warning-inline');
    expect(block.text()).toContain('No monitor');
  });
});
