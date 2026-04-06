import { flushPromises, mount } from "@vue/test-utils";
import InteractiveFloorPlan from "../InteractiveFloorPlan.vue";
import { ApiError } from "../../api/client";
import { createTestI18n } from "../../__tests__/helpers/i18n";
import { fetchFloorPlanPositions } from "../../api/floorPlanPositions";
import { fetchItems } from "../../api/items";
import {
  createBooking,
  createMultiDayBooking,
  cancelBooking,
} from "../../api/bookings";
import { fetchItemGroups } from "../../api/itemGroups";
import { fetchAreas } from "../../api/areas";

vi.mock("../../api/floorPlanPositions", () => ({
  fetchFloorPlanPositions: vi.fn(),
}));
vi.mock("../../api/items", () => ({ fetchItems: vi.fn() }));
vi.mock("../../api/bookings", () => ({
  createBooking: vi.fn(),
  createMultiDayBooking: vi.fn(),
  cancelBooking: vi.fn(),
}));
vi.mock("../../api/itemGroups", () => ({ fetchItemGroups: vi.fn() }));
vi.mock("../../api/areas", () => ({ fetchAreas: vi.fn() }));

describe("InteractiveFloorPlan", () => {
  const fetchFloorPlanPositionsMock = vi.mocked(fetchFloorPlanPositions);
  const fetchItemsMock = vi.mocked(fetchItems);
  const createBookingMock = vi.mocked(createBooking);
  const createMultiDayBookingMock = vi.mocked(createMultiDayBooking);
  const cancelBookingMock = vi.mocked(cancelBooking);
  const fetchItemGroupsMock = vi.mocked(fetchItemGroups);
  const fetchAreasMock = vi.mocked(fetchAreas);
  const originalMatchMedia = window.matchMedia;
  const originalInnerWidth = window.innerWidth;
  const originalInnerHeight = window.innerHeight;

  const weekDates = [
    "2026-04-06",
    "2026-04-07",
    "2026-04-08",
    "2026-04-09",
    "2026-04-10",
    "2026-04-11",
    "2026-04-12",
  ];

  const stubs = {
    "v-btn": {
      template: '<button type="button" v-bind="$attrs"><slot /></button>',
    },
    "v-tooltip": {
      template: '<div><slot name="activator" :props="{}" /><slot /></div>',
    },
    "v-icon": {
      template: "<span><slot /></span>",
    },
    "v-spacer": {
      template: "<div />",
    },
    "v-progress-circular": {
      template: "<div />",
    },
    "v-alert": {
      template: '<div v-bind="$attrs"><slot /></div>',
    },
    "v-snackbar": {
      props: ["modelValue"],
      template:
        '<div v-if="modelValue" v-bind="$attrs"><slot /><slot name="actions" /></div>',
    },
    "v-dialog": {
      props: ["modelValue", "fullscreen", "persistent"],
      template:
        '<div v-if="modelValue" v-bind="$attrs" :data-fullscreen="fullscreen" :data-persistent="persistent"><slot /></div>',
    },
    "v-card": {
      template: "<div><slot /></div>",
    },
    "v-card-title": {
      template: '<div v-bind="$attrs"><slot /></div>',
    },
    "v-card-text": {
      template: '<div v-bind="$attrs"><slot /></div>',
    },
    "v-card-actions": {
      template: "<div><slot /></div>",
    },
    "v-checkbox": {
      props: ["modelValue", "disabled", "color"],
      emits: ["update:modelValue"],
      template:
        '<div v-bind="$attrs" :data-disabled="disabled" :data-color="color"><input type="checkbox" :checked="modelValue" :disabled="disabled" @change="$emit(\'update:modelValue\', !modelValue)" /><slot /></div>',
    },
    "v-chip": {
      template: '<span v-bind="$attrs"><slot /></span>',
    },
  };

  function setViewport(width: number, height = 900) {
    Object.defineProperty(window, "innerWidth", {
      configurable: true,
      value: width,
    });
    Object.defineProperty(window, "innerHeight", {
      configurable: true,
      value: height,
    });
    window.matchMedia = vi.fn().mockImplementation((query: string) => ({
      matches:
        query === "(max-width: 768px)"
          ? width <= 768
          : query === "(max-height: 500px)"
            ? height <= 500
            : false,
      media: query,
      onchange: null,
      addListener: vi.fn(),
      removeListener: vi.fn(),
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
      dispatchEvent: vi.fn(),
    })) as typeof window.matchMedia;
  }

  function basePosition(itemId = "item-1", label = "Desk A") {
    return {
      id: `pos-${itemId}`,
      type: "floor-plan-positions",
      attributes: {
        item_id: itemId,
        label,
        x: 10,
        y: 20,
        width: 18,
        height: 12,
        border_width: 2,
      },
    };
  }

  function mountComponent(
    props: Partial<InstanceType<typeof InteractiveFloorPlan>["$props"]> = {},
  ) {
    return mount(InteractiveFloorPlan, {
      props: {
        floorPlan: "level-2.png",
        title: "Cube 1",
        weekLabel: "CW 15",
        weekDates,
        itemGroupId: "ig-1",
        ...props,
      },
      global: {
        plugins: [createTestI18n()],
        stubs,
      },
    });
  }

  beforeEach(() => {
    localStorage.clear();
    setViewport(1280, 900);
    fetchFloorPlanPositionsMock.mockResolvedValue({
      data: [basePosition()],
    } as never);
    fetchItemsMock.mockResolvedValue({
      data: [
        {
          id: "item-1",
          type: "items",
          attributes: {
            name: "Desk A",
            equipment: ["Monitor"],
            availability: "available",
            booked_by_me: false,
          },
        },
      ],
    } as never);
    createBookingMock.mockResolvedValue({ data: { id: "booking-1" } } as never);
    createMultiDayBookingMock.mockResolvedValue({
      created: [
        { id: "booking-1", type: "bookings", attributes: {} },
        { id: "booking-2", type: "bookings", attributes: {} },
      ],
      conflicts: [],
    } as never);
    cancelBookingMock.mockResolvedValue(undefined as never);
    fetchItemGroupsMock.mockResolvedValue({ data: [] } as never);
    fetchAreasMock.mockResolvedValue({ data: [] } as never);
  });

  afterEach(() => {
    window.matchMedia = originalMatchMedia;
    Object.defineProperty(window, "innerWidth", {
      configurable: true,
      value: originalInnerWidth,
    });
    Object.defineProperty(window, "innerHeight", {
      configurable: true,
      value: originalInnerHeight,
    });
  });

  it("opens a confirmation dialog before booking", async () => {
    const wrapper = mountComponent();
    await flushPromises();

    await wrapper.get('[data-cy="fp-item-item-1"]').trigger("click");
    await flushPromises();

    expect(createBookingMock).not.toHaveBeenCalled();
    expect(wrapper.find('[data-cy="fp-booking-dialog"]').exists()).toBe(true);
    expect(wrapper.get('[data-cy="fp-booking-summary"]').text()).toContain(
      "Book Desk A",
    );
  });

  it("creates a multi-day booking from the selected week", async () => {
    const wrapper = mountComponent();
    await flushPromises();

    await wrapper.get('[data-cy="fp-item-item-1"]').trigger("click");
    await flushPromises();
    await wrapper
      .get('[data-cy="fp-booking-day-TU"] input')
      .trigger("change");
    await wrapper.get('[data-cy="fp-booking-confirm"]').trigger("click");
    await flushPromises();

    expect(createMultiDayBookingMock).toHaveBeenCalledWith("item-1", [
      "2026-04-06",
      "2026-04-07",
    ]);
    expect(wrapper.find('[data-cy="fp-booking-success"]').exists()).toBe(true);
  });

  it("drills into the detailed floor plan when a child floor plan exists", async () => {
    fetchFloorPlanPositionsMock
      .mockResolvedValueOnce({
        data: [basePosition("ig-cube-2", "Cube 2")],
      } as never)
      .mockResolvedValueOnce({
        data: [basePosition("item-2", "Desk 07")],
      } as never);
    fetchAreasMock.mockResolvedValue({
      data: [{ id: "area-1", type: "areas", attributes: { name: "Office" } }],
    } as never);
    fetchItemGroupsMock.mockResolvedValue({
      data: [
        {
          id: "ig-cube-2",
          type: "item-groups",
          attributes: {
            name: "Cube 2",
            floor_plan: "cube-2.png",
          },
        },
      ],
    } as never);
    fetchItemsMock.mockResolvedValue({
      data: [
        {
          id: "item-2",
          type: "items",
          attributes: {
            name: "Desk 07",
            equipment: [],
            availability: "available",
            booked_by_me: false,
          },
        },
      ],
    } as never);

    const wrapper = mountComponent({
      floorPlan: "office.png",
      title: "Office",
      itemGroupId: "",
      areaLevel: true,
    });
    await flushPromises();

    await wrapper.get('[data-cy="fp-area-ig-cube-2"]').trigger("click");
    await flushPromises();

    expect(fetchFloorPlanPositionsMock).toHaveBeenLastCalledWith("cube-2.png");
    expect(wrapper.get('[data-cy="fp-breadcrumb-current"]').text()).toBe(
      "Cube 2",
    );
  });

  it("disables days where the selected item is already booked", async () => {
    fetchItemsMock.mockImplementation(
      async (_itemGroupId: string, date: string) =>
        ({
          data: [
            {
              id: "item-1",
              type: "items",
              attributes: {
                name: "Desk A",
                equipment: [],
                availability: date === "2026-04-08" ? "occupied" : "available",
                booked_by_me: false,
                booker_name: date === "2026-04-08" ? "Jane Doe" : undefined,
              },
            },
          ],
        }) as never,
    );

    const wrapper = mountComponent();
    await flushPromises();

    await wrapper.get('[data-cy="fp-item-item-1"]').trigger("click");
    await flushPromises();

    const weRow = wrapper.get('[data-cy="fp-booking-day-WE"]');
    const checkbox = weRow.find("input");
    expect((checkbox.element as HTMLInputElement).disabled).toBe(true);
    expect(weRow.text()).toContain("Jane Doe");
  });

  it("shows all weekdays including weekends on compact screens and highlights the selection in red", async () => {
    setViewport(430, 900);

    const wrapper = mountComponent();
    await flushPromises();

    expect(wrapper.find('[data-cy="fp-day-SA"]').exists()).toBe(true);
    expect(wrapper.find('[data-cy="fp-day-SU"]').exists()).toBe(true);
    const moBtn = wrapper.get('[data-cy="fp-day-MO"]');
    expect(moBtn.attributes("variant")).toBe("flat");
    expect(moBtn.attributes("color")).toBe("error");
  });

  it("renders busy-item avatars and hides them when show avatars is turned off", async () => {
    fetchItemsMock.mockResolvedValue({
      data: [
        {
          id: "item-1",
          type: "items",
          attributes: {
            name: "Desk A",
            equipment: ["Monitor"],
            availability: "occupied",
            booker_name: "Alice Smith",
            booker_user_id: "user-1",
            booked_by_me: false,
          },
        },
      ],
    } as never);

    const wrapper = mountComponent();
    await flushPromises();

    const avatar = wrapper.get('[data-cy="fp-avatar-item-1"]');
    expect(avatar.attributes("src")).toBe("/api/v1/avatars/user-1");

    await wrapper.get('[data-cy="fp-show-avatars"] input').trigger("change");
    await flushPromises();

    expect(localStorage.getItem("sithub_fp_show_avatars")).toBe("false");
    expect(wrapper.find('[data-cy="fp-avatar-item-1"]').exists()).toBe(false);
  });

  it("renders reserved items with a lock and blocks booking interaction", async () => {
    fetchItemsMock.mockResolvedValue({
      data: [
        {
          id: "item-1",
          type: "items",
          attributes: {
            name: "Desk A",
            equipment: ["Monitor"],
            availability: "available",
            reserved: true,
            booked_by_me: false,
          },
        },
      ],
    } as never);

    const wrapper = mountComponent();
    await flushPromises();

    expect(wrapper.find('[data-cy="fp-lock-item-1"]').exists()).toBe(true);

    await wrapper.get('[data-cy="fp-item-item-1"]').trigger("click");
    await flushPromises();

    expect(wrapper.find('[data-cy="fp-booking-dialog"]').exists()).toBe(false);
    expect(createBookingMock).not.toHaveBeenCalled();
  });

  it("uses close as back when drilled into a detailed floor plan", async () => {
    fetchFloorPlanPositionsMock
      .mockResolvedValueOnce({
        data: [basePosition("ig-cube-2", "Cube 2")],
      } as never)
      .mockResolvedValueOnce({
        data: [basePosition("item-2", "Desk 07")],
      } as never)
      .mockResolvedValueOnce({
        data: [basePosition("ig-cube-2", "Cube 2")],
      } as never);
    fetchAreasMock.mockResolvedValue({
      data: [{ id: "area-1", type: "areas", attributes: { name: "Office" } }],
    } as never);
    fetchItemGroupsMock.mockResolvedValue({
      data: [
        {
          id: "ig-cube-2",
          type: "item-groups",
          attributes: {
            name: "Cube 2",
            floor_plan: "cube-2.png",
          },
        },
      ],
    } as never);
    fetchItemsMock.mockResolvedValue({
      data: [
        {
          id: "item-2",
          type: "items",
          attributes: {
            name: "Desk 07",
            equipment: [],
            availability: "available",
            booked_by_me: false,
          },
        },
      ],
    } as never);

    const wrapper = mountComponent({
      floorPlan: "office.png",
      title: "Office",
      itemGroupId: "",
      areaLevel: true,
    });
    await flushPromises();

    await wrapper.get('[data-cy="fp-area-ig-cube-2"]').trigger("click");
    await flushPromises();
    await wrapper.get('[data-cy="fp-close-btn"]').trigger("click");
    await flushPromises();

    expect(wrapper.emitted("close")).toBeFalsy();
    expect(wrapper.find('[data-cy="fp-breadcrumb-current"]').exists()).toBe(false);
    expect(fetchFloorPlanPositionsMock).toHaveBeenLastCalledWith("office.png");
  });

  it("shows a precise conflict error when booking fails", async () => {
    createBookingMock.mockRejectedValue(
      new ApiError("Request failed: 409", 409, "Item is already booked for this date"),
    );

    const wrapper = mountComponent();
    await flushPromises();

    await wrapper.get('[data-cy="fp-item-item-1"]').trigger("click");
    await flushPromises();
    await wrapper.get('[data-cy="fp-booking-confirm"]').trigger("click");
    await flushPromises();

    expect(wrapper.get('[data-cy="fp-booking-error"]').text()).toContain(
      "The selected item is already booked.",
    );
  });

  it("preserves the mobile auto-fit zoom when only the week changes", async () => {
    setViewport(430, 900);

    const wrapper = mountComponent();
    await flushPromises();

    const image = wrapper.get(".fp-image-fit");
    const fitContainer = image.element.parentElement?.parentElement as HTMLElement | null;
    expect(fitContainer).not.toBeNull();

    Object.defineProperty(image.element, "naturalWidth", {
      configurable: true,
      value: 1200,
    });
    Object.defineProperty(fitContainer, "clientWidth", {
      configurable: true,
      value: 600,
    });

    await image.trigger("load");
    await flushPromises();
    expect(wrapper.text()).toContain("75%");

    await wrapper.setProps({
      weekDates: [
        "2026-04-13",
        "2026-04-14",
        "2026-04-15",
        "2026-04-16",
        "2026-04-17",
        "2026-04-18",
        "2026-04-19",
      ],
    });
    await flushPromises();

    expect(wrapper.text()).toContain("75%");
  });
});
