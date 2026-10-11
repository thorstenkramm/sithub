import { mount } from '@vue/test-utils';
import PageHeader from '../PageHeader.vue';

describe('PageHeader', () => {
  it('renders title', () => {
    const wrapper = mount(PageHeader, {
      props: { title: 'Test Title' },
      global: {
        stubs: {
          'v-icon': true,
          'router-link': true
        }
      }
    });
    expect(wrapper.text()).toContain('Test Title');
  });

  it('renders subtitle when provided', () => {
    const wrapper = mount(PageHeader, {
      props: { title: 'Title', subtitle: 'Subtitle text' },
      global: {
        stubs: {
          'v-icon': true,
          'router-link': true
        }
      }
    });
    expect(wrapper.text()).toContain('Subtitle text');
  });

  it('renders breadcrumbs when provided', () => {
    const breadcrumbs = [
      { text: 'Home', to: '/' },
      { text: 'Areas', to: '/areas' },
      { text: 'Current' }
    ];
    const wrapper = mount(PageHeader, {
      props: { title: 'Title', breadcrumbs },
      global: {
        stubs: {
          'v-icon': true,
          'router-link': {
            template: '<a><slot /></a>',
            props: ['to']
          }
        }
      }
    });
    expect(wrapper.text()).toContain('Home');
    expect(wrapper.text()).toContain('Areas');
    expect(wrapper.text()).toContain('Current');
  });

  it('renders actions slot', () => {
    const wrapper = mount(PageHeader, {
      props: { title: 'Title' },
      slots: {
        actions: '<button>Action</button>'
      },
      global: {
        stubs: {
          'v-icon': true,
          'router-link': true
        }
      }
    });
    expect(wrapper.find('button').exists()).toBe(true);
  });

  it('should keep crumb entries whole and let the trail wrap between entries', () => {
    const breadcrumbs = [
      { text: 'Startseite', to: '/' },
      { text: 'Büro 2.EG', to: '/areas/buero_2eg/item-groups' },
      { text: 'Cube 1' }
    ];
    const wrapper = mount(PageHeader, {
      props: { title: 'Title', breadcrumbs },
      global: {
        stubs: {
          'v-icon': { template: '<i><slot /></i>' },
          'router-link': { props: ['to'], template: '<a v-bind="$attrs"><slot /></a>' }
        }
      }
    });

    expect(wrapper.get('[data-cy="breadcrumbs"] ol').classes()).toContain('flex-wrap');

    const link = wrapper.get('[data-cy="breadcrumb-link-1"]');
    expect(link.classes()).toContain('breadcrumb-label');
    expect(link.attributes('title')).toBe('Büro 2.EG');

    const current = wrapper.get('[data-cy="breadcrumb-text-2"]');
    expect(current.classes()).toContain('breadcrumb-label');
    expect(current.attributes('title')).toBe('Cube 1');
  });
});
