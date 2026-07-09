import { describe, it, expect, vi, beforeEach } from 'vitest';
import { mount, flushPromises } from '@vue/test-utils';
import ColleagueSelect from '../ColleagueSelect.vue';
import { fetchColleagues } from '../../api/users';
import { createTestI18n } from '../../__tests__/helpers/i18n';

vi.mock('../../api/users', () => ({ fetchColleagues: vi.fn() }));
const fetchColleaguesMock = vi.mocked(fetchColleagues);

const stubs = {
  'v-radio-group': {
    props: ['modelValue'],
    emits: ['update:modelValue'],
    template: '<div v-bind="$attrs"><slot /></div>'
  },
  'v-radio': {
    props: ['label', 'value'],
    template: '<label v-bind="$attrs"><slot /></label>'
  },
  'v-expand-transition': { template: '<div><slot /></div>' },
  'v-autocomplete': {
    props: ['modelValue', 'items'],
    emits: ['update:modelValue'],
    template: '<select v-bind="$attrs"></select>'
  }
};

function mountSelect(prefix = 'test', modelValue: string | null = null) {
  return mount(ColleagueSelect, {
    props: { dataCyPrefix: prefix, modelValue },
    global: { stubs, plugins: [createTestI18n()] }
  });
}

describe('ColleagueSelect', () => {
  beforeEach(() => {
    fetchColleaguesMock.mockReset();
    fetchColleaguesMock.mockResolvedValue({
      data: [{ id: 'u-1', type: 'colleagues', attributes: { display_name: 'Jane Doe' } }]
    } as never);
  });

  it('renders self/colleague radios with the given data-cy prefix', async () => {
    const wrapper = mountSelect('fp');
    await flushPromises();

    expect(wrapper.find('[data-cy="fp-colleague-select"]').exists()).toBe(true);
    expect(wrapper.find('[data-cy="fp-book-self-radio"]').exists()).toBe(true);
    expect(wrapper.find('[data-cy="fp-book-colleague-radio"]').exists()).toBe(true);
    // Autocomplete hidden while in "self" mode.
    expect(wrapper.find('[data-cy="fp-colleague-autocomplete"]').exists()).toBe(false);
  });

  it('starts in colleague mode and loads colleagues when a model value is preset', async () => {
    const wrapper = mountSelect('tile', 'u-1');
    await flushPromises();

    expect(fetchColleaguesMock).toHaveBeenCalled();
    expect(wrapper.find('[data-cy="tile-colleague-autocomplete"]').exists()).toBe(true);
  });
});
