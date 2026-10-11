import { describe, it, expect } from 'vitest';
import { mount } from '@vue/test-utils';
import WeekSelector from '../WeekSelector.vue';
import { createTestI18n } from '../../__tests__/helpers/i18n';

const options = [
  { label: '05.10.-11.10.2026 - Week 41', value: '2026-W41' },
  { label: '12.10.-18.10.2026 - Week 42', value: '2026-W42' },
  { label: '19.10.-25.10.2026 - Week 43', value: '2026-W43' }
];

const selectStub = {
  name: 'VSelectStub',
  props: ['modelValue', 'items', 'label'],
  emits: ['update:modelValue'],
  template: '<select v-bind="$attrs" :data-value="modelValue" :data-label="label"></select>'
};

const stubs = {
  'v-select': selectStub,
  'v-btn': {
    props: ['disabled'],
    template: '<button v-bind="$attrs" :disabled="disabled"><slot /></button>'
  },
  'v-icon': { template: '<i><slot /></i>' }
};

function mountSelector(modelValue: string, selectDataCy?: string) {
  return mount(WeekSelector, {
    props: { modelValue, options, label: 'Calendar Week', ...(selectDataCy ? { selectDataCy } : {}) },
    global: { stubs, plugins: [createTestI18n()] }
  });
}

describe('WeekSelector', () => {
  it('should render the select with the default data-cy and the label', () => {
    const wrapper = mountSelector('2026-W41');

    const select = wrapper.get('[data-cy="week-selector"]');
    expect(select.attributes('data-value')).toBe('2026-W41');
    expect(select.attributes('data-label')).toBe('Calendar Week');
  });

  it('should use a custom data-cy for the select', () => {
    const wrapper = mountSelector('2026-W41', 'fp-week-selector');

    expect(wrapper.find('[data-cy="fp-week-selector"]').exists()).toBe(true);
    expect(wrapper.find('[data-cy="week-selector"]').exists()).toBe(false);
  });

  it('should disable previous week on the first option', () => {
    const wrapper = mountSelector('2026-W41');

    expect(wrapper.get('[data-cy="week-prev"]').attributes('disabled')).toBeDefined();
    expect(wrapper.get('[data-cy="week-next"]').attributes('disabled')).toBeUndefined();
  });

  it('should disable next week on the last option', () => {
    const wrapper = mountSelector('2026-W43');

    expect(wrapper.get('[data-cy="week-prev"]').attributes('disabled')).toBeUndefined();
    expect(wrapper.get('[data-cy="week-next"]').attributes('disabled')).toBeDefined();
  });

  it('should disable both arrows when the value is not an option', () => {
    const wrapper = mountSelector('2027-W01');

    expect(wrapper.get('[data-cy="week-prev"]').attributes('disabled')).toBeDefined();
    expect(wrapper.get('[data-cy="week-next"]').attributes('disabled')).toBeDefined();
  });

  it('should step forward by exactly one week', async () => {
    const wrapper = mountSelector('2026-W41');

    await wrapper.get('[data-cy="week-next"]').trigger('click');

    expect(wrapper.emitted('update:modelValue')).toEqual([['2026-W42']]);
  });

  it('should step back by exactly one week', async () => {
    const wrapper = mountSelector('2026-W43');

    await wrapper.get('[data-cy="week-prev"]').trigger('click');

    expect(wrapper.emitted('update:modelValue')).toEqual([['2026-W42']]);
  });

  it('should forward a selection from the dropdown', async () => {
    const wrapper = mountSelector('2026-W41');

    await wrapper.getComponent(selectStub).vm.$emit('update:modelValue', '2026-W43');

    expect(wrapper.emitted('update:modelValue')).toEqual([['2026-W43']]);
  });

  it('should have accessible labels on the arrows', () => {
    const wrapper = mountSelector('2026-W42');

    expect(wrapper.get('[data-cy="week-prev"]').attributes('aria-label')).toBe('Previous week');
    expect(wrapper.get('[data-cy="week-next"]').attributes('aria-label')).toBe('Next week');
  });
});
