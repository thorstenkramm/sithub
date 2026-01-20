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
});
