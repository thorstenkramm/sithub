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
  fetchMyBookings,
} from "../../api/bookings";
import { fetchColleagues } from "../../api/users";
import { fetchItemGroups } from "../../api/itemGroups";
import { fetchAreas } from "../../api/areas";
import { __resetLegacyPurgeForTests } from "../../composables/useFavorites";

const liveFeed = vi.hoisted(() => ({
  handler: null as ((event: unknown) => void) | null,
}));

vi.mock("../../api/floorPlanPositions", () => ({
  fetchFloorPlanPositions: vi.fn(),
}));
vi.mock("../../api/items", () => ({ fetchItems: vi.fn() }));
vi.mock("../../api/bookings", () => ({
  createBooking: vi.fn(),
  createMultiDayBooking: vi.fn(),
  cancelBooking: vi.fn(),
  fetchMyBookings: vi.fn(),
}));
vi.mock("../../api/users", () => ({ fetchColleagues: vi.fn() }));
vi.mock("../../api/itemGroups", () => ({ fetchItemGroups: vi.fn() }));
vi.mock("../../api/areas", () => ({ fetchAreas: vi.fn() }));
vi.mock("../../stores/useLiveFeedStore", () => ({
  useLiveFeedStore: () => ({
    start: vi.fn(),
    stop: vi.fn(),
    reset: vi.fn(),
    subscribe: (handler: (event: unknown) => void) => {
      liveFeed.handler = handler;
      return () => {
        if (liveFeed.handler === handler) {
          liveFeed.handler = null;
        }
      };
    },
  }),
}));

describe("InteractiveFloorPlan", () => {
  const fetchFloorPlanPositionsMock = vi.mocked(fetchFloorPlanPositions);
  const fetchItemsMock = vi.mocked(fetchItems);
  const createBookingMock = vi.mocked(createBooking);
  const createMultiDayBookingMock = vi.mocked(createMultiDayBooking);
  const cancelBookingMock = vi.mocked(cancelBooking);
  const fetchMyBookingsMock = vi.mocked(fetchMyBookings);
  const fetchColleaguesMock = vi.mocked(fetchColleagues);
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
    "v-radio-group": {
      props: ["modelValue"],
      template: '<div v-bind="$attrs"><slot /></div>',
    },
    "v-radio": {
      template: '<label v-bind="$attrs"><slot /></label>',
    },
    "v-autocomplete": {
      props: ["modelValue"],
      template: '<div v-bind="$attrs"><slot /></div>',
    },
    "v-expand-transition": {
      template: "<div><slot /></div>",
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

  function mockAreaLevelDeskScenario() {
    fetchFloorPlanPositionsMock.mockImplementation(async (floorPlan: string) => ({
      data:
        floorPlan === "office.png"
          ? [basePosition("ig-cube-2", "Cube 2"), basePosition("item-2", "Desk 07")]
          : [basePosition("item-2", "Desk 07")],
    }) as never);
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
    fetchItemsMock.mockImplementation(async () => ({
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
    }) as never);
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
    vi.useFakeTimers();
    vi.setSystemTime(new Date("2026-04-08T10:00:00"));
    localStorage.clear();
    __resetLegacyPurgeForTests();
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
    fetchMyBookingsMock.mockResolvedValue({ data: [] } as never);
    fetchColleaguesMock.mockResolvedValue({ data: [] } as never);
    fetchItemGroupsMock.mockResolvedValue({ data: [] } as never);
    fetchAreasMock.mockResolvedValue({ data: [] } as never);
    liveFeed.handler = null;
  });

  afterEach(() => {
    vi.useRealTimers();
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

  it("shows the warning confirmation before booking a warned free item, then books on confirm", async () => {
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
            warning: "Standing desk",
          },
        },
      ],
    } as never);
    const wrapper = mountComponent();
    await flushPromises();

    await wrapper.get('[data-cy="fp-item-item-1"]').trigger("click");
    await flushPromises();

    // Confirming the date selection must surface the shared warning dialog
    // before any booking is created.
    await wrapper.get('[data-cy="fp-booking-confirm"]').trigger("click");
    await flushPromises();

    expect(wrapper.find('[data-cy="warning-dialog"]').exists()).toBe(true);
    expect(wrapper.find('[data-cy="warning-message"]').text()).toContain(
      "Standing desk",
    );
    expect(createBookingMock).not.toHaveBeenCalled();

    // Confirming the warning proceeds with the single-day booking.
    await wrapper.get('[data-cy="warning-confirm-btn"]').trigger("click");
    await flushPromises();
    expect(createBookingMock).toHaveBeenCalledWith("item-1", "2026-04-08", undefined);
  });

  it("aborts the floor-plan booking when the warning confirmation is cancelled", async () => {
    // The suite's beforeEach re-sets the resolved value but does not clear call
    // history, so reset it explicitly for the not-called assertion below.
    createBookingMock.mockClear();
    fetchItemsMock.mockResolvedValue({
      data: [
        {
          id: "item-1",
          type: "items",
          attributes: {
            name: "Desk A",
            equipment: [],
            availability: "available",
            booked_by_me: false,
            warning: "Standing desk",
          },
        },
      ],
    } as never);
    const wrapper = mountComponent();
    await flushPromises();

    await wrapper.get('[data-cy="fp-item-item-1"]').trigger("click");
    await flushPromises();
    await wrapper.get('[data-cy="fp-booking-confirm"]').trigger("click");
    await flushPromises();
    expect(wrapper.find('[data-cy="warning-dialog"]').exists()).toBe(true);

    await wrapper.get('[data-cy="warning-cancel-btn"]').trigger("click");
    await flushPromises();

    expect(createBookingMock).not.toHaveBeenCalled();
    expect(wrapper.find('[data-cy="warning-dialog"]').exists()).toBe(false);
  });

  it("creates a multi-day booking from the selected week", async () => {
    const wrapper = mountComponent();
    await flushPromises();

    await wrapper.get('[data-cy="fp-item-item-1"]').trigger("click");
    await flushPromises();

    // Select an additional future day (today is pre-selected via initializeBookingSelection)
    const today = new Date().toISOString().slice(0, 10);
    const todayIdx = weekDates.indexOf(today);
    // Find the next non-past day after today to toggle on
    const nextIdx = todayIdx + 1;
    const nextDate = weekDates[nextIdx];
    const nextDayIndex = new Date(nextDate!).getDay();
    const labels = ["SU", "MO", "TU", "WE", "TH", "FR", "SA"];
    const nextLabel = labels[nextDayIndex]!;

    await wrapper
      .get(`[data-cy="fp-booking-day-${nextLabel}"] input`)
      .trigger("change");
    await wrapper.get('[data-cy="fp-booking-confirm"]').trigger("click");
    await flushPromises();

    expect(createMultiDayBookingMock).toHaveBeenCalledWith(
      "item-1",
      [today, nextDate],
      undefined,
    );
    expect(wrapper.find('[data-cy="fp-booking-success"]').exists()).toBe(true);
  });

  it("refreshes the visible item group when a relevant live event arrives", async () => {
    mountComponent();
    await flushPromises();

    fetchItemsMock.mockClear();
    expect(liveFeed.handler).toBeTypeOf("function");
    liveFeed.handler!({
      type: "booking.created",
      booking_id: "booking-1",
      item_id: "item-1",
      user_id: "other-user",
      booking_date: "2026-04-08",
      timestamp: "2026-05-10T12:00:00Z",
    });

    await vi.advanceTimersByTimeAsync(300);
    await flushPromises();

    expect(fetchItemsMock).toHaveBeenCalledTimes(1);
    expect(fetchItemsMock).toHaveBeenCalledWith("ig-1", "2026-04-08");
  });

  it("renders a removable favorite heart on free floor-plan items", async () => {
    localStorage.setItem("sithub_favorite_items", JSON.stringify([{
      areaId: "area-1",
      itemId: "item-1",
      itemName: "Desk A",
      itemGroupId: "ig-1",
      itemGroupName: "Cube 1",
    }]));

    const wrapper = mountComponent({ areaId: "area-1" });
    await flushPromises();

    const heart = wrapper.get('[data-cy="fp-favorite-heart-item-1"]');
    await heart.trigger("click");
    await flushPromises();

    expect(wrapper.find('[data-cy="fp-favorite-heart-item-1"]').exists()).toBe(false);
    expect(localStorage.getItem("sithub_favorite_items")).toBe("[]");
    expect(wrapper.find('[data-cy="fp-booking-dialog"]').exists()).toBe(false);
  });

  it("does not render a favorite heart on busy floor-plan items", async () => {
    localStorage.setItem("sithub_favorite_items", JSON.stringify([{
      areaId: "area-1",
      itemId: "item-1",
      itemName: "Desk A",
      itemGroupId: "ig-1",
      itemGroupName: "Cube 1",
    }]));
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

    const wrapper = mountComponent({ areaId: "area-1" });
    await flushPromises();

    expect(wrapper.find('[data-cy="fp-favorite-heart-item-1"]').exists()).toBe(false);
    expect(wrapper.find('[data-cy="fp-avatar-item-1"]').exists()).toBe(true);
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

  it("books desks directly by default on large screens", async () => {
    mockAreaLevelDeskScenario();

    const wrapper = mountComponent({
      floorPlan: "office.png",
      title: "Office",
      itemGroupId: "",
      areaLevel: true,
    });
    await flushPromises();

    const drillDownToggle = wrapper.get(
      '[data-cy="floor-plan-area-drill-down-toggle"] input',
    );
    expect((drillDownToggle.element as HTMLInputElement).checked).toBe(false);

    fetchFloorPlanPositionsMock.mockClear();
    await wrapper.get('[data-cy="fp-desk-item-2"]').trigger("click");
    await flushPromises();

    expect(wrapper.find('[data-cy="fp-booking-dialog"]').exists()).toBe(true);
    expect(wrapper.find('[data-cy="fp-breadcrumb-current"]').exists()).toBe(
      false,
    );
    expect(fetchFloorPlanPositionsMock).not.toHaveBeenCalledWith("cube-2.png");
  });

  it("drills into desk clicks by default on compact screens", async () => {
    mockAreaLevelDeskScenario();
    setViewport(430, 900);

    const wrapper = mountComponent({
      floorPlan: "office.png",
      title: "Office",
      itemGroupId: "",
      areaLevel: true,
    });
    await flushPromises();

    const drillDownToggle = wrapper.get(
      '[data-cy="floor-plan-area-drill-down-toggle"] input',
    );
    expect((drillDownToggle.element as HTMLInputElement).checked).toBe(true);

    await wrapper.get('[data-cy="fp-desk-item-2"]').trigger("click");
    await flushPromises();

    expect(wrapper.find('[data-cy="fp-booking-dialog"]').exists()).toBe(false);
    expect(fetchFloorPlanPositionsMock).toHaveBeenLastCalledWith("cube-2.png");
    expect(wrapper.get('[data-cy="fp-breadcrumb-current"]').text()).toBe(
      "Cube 2",
    );
  });

  it("re-applies the viewport default on resize until the user chooses a value", async () => {
    mockAreaLevelDeskScenario();

    const wrapper = mountComponent({
      floorPlan: "office.png",
      title: "Office",
      itemGroupId: "",
      areaLevel: true,
    });
    await flushPromises();

    const drillDownToggleSelector =
      '[data-cy="floor-plan-area-drill-down-toggle"] input';
    expect(
      (wrapper.get(drillDownToggleSelector).element as HTMLInputElement)
        .checked,
    ).toBe(false);

    setViewport(430, 900);
    window.dispatchEvent(new Event("resize"));
    await flushPromises();
    expect(
      (wrapper.get(drillDownToggleSelector).element as HTMLInputElement)
        .checked,
    ).toBe(true);

    setViewport(1280, 900);
    window.dispatchEvent(new Event("resize"));
    await flushPromises();
    expect(
      (wrapper.get(drillDownToggleSelector).element as HTMLInputElement)
        .checked,
    ).toBe(false);
  });

  it("persists the drill-down toggle override across floor plan sessions", async () => {
    mockAreaLevelDeskScenario();

    const props = {
      floorPlan: "office.png",
      title: "Office",
      itemGroupId: "",
      areaLevel: true,
    } as const;

    const firstWrapper = mountComponent(props);
    await flushPromises();

    await firstWrapper
      .get('[data-cy="floor-plan-area-drill-down-toggle"] input')
      .trigger("change");
    await flushPromises();

    expect(localStorage.getItem("sithub_area_drill_down")).toBe("on");
    firstWrapper.unmount();

    const secondWrapper = mountComponent(props);
    await flushPromises();

    expect(
      (
        secondWrapper.get(
          '[data-cy="floor-plan-area-drill-down-toggle"] input',
        ).element as HTMLInputElement
      ).checked,
    ).toBe(true);

    await secondWrapper.get('[data-cy="fp-desk-item-2"]').trigger("click");
    await flushPromises();

    expect(secondWrapper.find('[data-cy="fp-booking-dialog"]').exists()).toBe(
      false,
    );
    expect(fetchFloorPlanPositionsMock).toHaveBeenLastCalledWith("cube-2.png");
  });

  it("disables days where the selected item is already booked", async () => {
    // Use a future date that's in the week range as the occupied day
    const today = new Date().toISOString().slice(0, 10);
    const todayIdx = weekDates.indexOf(today);
    const futureIdx = todayIdx + 1;
    const occupiedDate = weekDates[futureIdx]!;
    const occupiedDayIndex = new Date(occupiedDate).getDay();
    const labels = ["SU", "MO", "TU", "WE", "TH", "FR", "SA"];
    const occupiedLabel = labels[occupiedDayIndex]!;

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
                availability: date === occupiedDate ? "occupied" : "available",
                booked_by_me: false,
                booker_name: date === occupiedDate ? "Jane Doe" : undefined,
              },
            },
          ],
        }) as never,
    );

    const wrapper = mountComponent();
    await flushPromises();

    await wrapper.get('[data-cy="fp-item-item-1"]').trigger("click");
    await flushPromises();

    const dayRow = wrapper.get(`[data-cy="fp-booking-day-${occupiedLabel}"]`);
    const checkbox = dayRow.find("input");
    expect((checkbox.element as HTMLInputElement).disabled).toBe(true);
    expect(dayRow.text()).toContain("Jane Doe");
  });

  it("shows all weekdays including weekends on compact screens and highlights the selection in red", async () => {
    setViewport(430, 900);

    const wrapper = mountComponent();
    await flushPromises();

    expect(wrapper.find('[data-cy="fp-day-SA"]').exists()).toBe(true);
    expect(wrapper.find('[data-cy="fp-day-SU"]').exists()).toBe(true);

    // The selected day defaults to today (WE) — verify it is highlighted
    const today = new Date().toISOString().slice(0, 10);
    const todayDayIndex = new Date(today).getDay(); // 0=SU,1=MO,...
    const labels = ["SU", "MO", "TU", "WE", "TH", "FR", "SA"];
    const todayLabel = labels[todayDayIndex] || "MO";
    const todayBtn = wrapper.get(`[data-cy="fp-day-${todayLabel}"]`);
    expect(todayBtn.attributes("variant")).toBe("flat");
    expect(todayBtn.attributes("color")).toBe("error");
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
    createBookingMock.mockClear();

    expect(wrapper.find('[data-cy="fp-lock-item-1"]').exists()).toBe(true);

    await wrapper.get('[data-cy="fp-item-item-1"]').trigger("click");
    await flushPromises();

    expect(wrapper.find('[data-cy="fp-booking-dialog"]').exists()).toBe(false);
    expect(createBookingMock).not.toHaveBeenCalled();
  });

  it("shows avatars for occupied reserved desks after drilling into a reserved area", async () => {
    fetchFloorPlanPositionsMock
      .mockResolvedValueOnce({
        data: [basePosition("ig-finance", "Finance")],
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
          id: "ig-finance",
          type: "item-groups",
          attributes: {
            name: "Finance",
            floor_plan: "finance.png",
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
            availability: "occupied",
            reserved: true,
            booker_name: "Alice Smith",
            booker_user_id: "user-1",
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

    await wrapper.get('[data-cy="fp-area-ig-finance"]').trigger("click");
    await flushPromises();

    expect(wrapper.find('[data-cy="fp-lock-item-2"]').exists()).toBe(false);
    expect(wrapper.get('[data-cy="fp-avatar-item-2"]').attributes("src")).toBe(
      "/api/v1/avatars/user-1",
    );
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

  it("preselects the selectedDay prop when it is a non-past day in the week (story 36.6)", async () => {
    // 2026-04-10 is a Friday, in the future relative to the 04-08 system time.
    const wrapper = mountComponent({ selectedDay: "2026-04-10" });
    await flushPromises();

    const frBtn = wrapper.get('[data-cy="fp-day-FR"]');
    expect(frBtn.attributes("variant")).toBe("flat");
    expect(frBtn.attributes("color")).toBe("error");
    // Availability is loaded for the preselected day, not today.
    expect(fetchItemsMock).toHaveBeenCalledWith("ig-1", "2026-04-10");
  });

  it("falls back to today when selectedDay is a past day (story 36.6)", async () => {
    // 2026-04-06 (Monday) is before the 04-08 system time.
    const wrapper = mountComponent({ selectedDay: "2026-04-06" });
    await flushPromises();

    // Today is Wednesday (WE) — it must be highlighted, not the past Monday.
    const weBtn = wrapper.get('[data-cy="fp-day-WE"]');
    expect(weBtn.attributes("variant")).toBe("flat");
    expect(weBtn.attributes("color")).toBe("error");
  });

  it("books every selected day on a colleague's behalf from the floor plan (story 36.7)", async () => {
    fetchColleaguesMock.mockResolvedValue({
      data: [{ id: "u-1", type: "colleagues", attributes: { display_name: "Jane Doe" } }],
    } as never);
    const wrapper = mountComponent();
    await flushPromises();

    await wrapper.get('[data-cy="fp-item-item-1"]').trigger("click");
    await flushPromises();

    // Select a second day so the multi-day path is used.
    const today = new Date().toISOString().slice(0, 10);
    const todayIdx = weekDates.indexOf(today);
    const nextDate = weekDates[todayIdx + 1]!;
    const nextLabel = ["SU", "MO", "TU", "WE", "TH", "FR", "SA"][
      new Date(nextDate).getDay()
    ]!;
    await wrapper
      .get(`[data-cy="fp-booking-day-${nextLabel}"] input`)
      .trigger("change");

    // Pick a colleague via the component model.
    (wrapper.vm as unknown as { bookingColleagueId: string | null }).bookingColleagueId = "u-1";
    await flushPromises();

    await wrapper.get('[data-cy="fp-booking-confirm"]').trigger("click");
    await flushPromises();

    expect(createMultiDayBookingMock).toHaveBeenCalledWith(
      "item-1",
      [today, nextDate],
      { forUserId: "u-1", forUserName: "Jane Doe" },
    );
  });

  it("prompts to swap an existing same area/day booking before booking (story 36.9)", async () => {
    fetchMyBookingsMock.mockResolvedValue({
      data: [
        {
          id: "existing-1",
          type: "bookings",
          attributes: {
            item_id: "other-desk",
            item_name: "Desk Z",
            area_id: "area-1",
            booking_date: "2026-04-08",
            booking_date_display: "",
            item_group_id: "ig-1",
            item_group_name: "Cube 1",
            area_name: "Office",
            created_at: "",
            booked_by_user_id: "me",
            booked_by_user_name: "Me",
            booked_for_me: true,
            note: "",
          },
        },
      ],
    } as never);

    const wrapper = mountComponent({ areaId: "area-1" });
    await flushPromises();
    createBookingMock.mockClear();

    await wrapper.get('[data-cy="fp-item-item-1"]').trigger("click");
    await flushPromises();
    await wrapper.get('[data-cy="fp-booking-confirm"]').trigger("click");
    await flushPromises();

    // Swap dialog appears; no booking yet.
    expect(wrapper.find('[data-cy="confirm-dialog-confirm"]').exists()).toBe(true);
    expect(createBookingMock).not.toHaveBeenCalled();

    // Confirm swaps: create new FIRST, then cancel old (story 36.9 D2).
    await wrapper.get('[data-cy="confirm-dialog-confirm"]').trigger("click");
    await flushPromises();
    expect(createBookingMock).toHaveBeenCalledWith("item-1", "2026-04-08", undefined);
    expect(cancelBookingMock).toHaveBeenCalledWith("existing-1");
    const createOrder = createBookingMock.mock.invocationCallOrder[0]!;
    const cancelOrder = cancelBookingMock.mock.invocationCallOrder[0]!;
    expect(createOrder).toBeLessThan(cancelOrder);
  });

  it("keeps the new booking and does not throw when the post-create cancel fails (story 36.9)", async () => {
    fetchMyBookingsMock.mockResolvedValue({
      data: [
        {
          id: "existing-1",
          type: "bookings",
          attributes: {
            item_id: "other-desk",
            item_name: "Desk Z",
            area_id: "area-1",
            booking_date: "2026-04-08",
            item_group_id: "ig-1",
            item_group_name: "Cube 1",
            area_name: "Office",
            created_at: "",
            booked_by_user_id: "me",
            booked_by_user_name: "Me",
            booked_for_me: true,
            note: "",
          },
        },
      ],
    } as never);

    const wrapper = mountComponent({ areaId: "area-1" });
    await flushPromises();
    createBookingMock.mockClear();
    cancelBookingMock.mockClear();
    createBookingMock.mockResolvedValue({ data: { id: "new-1" } } as never);
    cancelBookingMock.mockRejectedValueOnce(new Error("cancel failed"));

    await wrapper.get('[data-cy="fp-item-item-1"]').trigger("click");
    await flushPromises();
    await wrapper.get('[data-cy="fp-booking-confirm"]').trigger("click");
    await flushPromises();
    await wrapper.get('[data-cy="confirm-dialog-confirm"]').trigger("click");
    await flushPromises();

    // New booking kept; success snackbar shown despite the cancel failure.
    expect(createBookingMock).toHaveBeenCalledWith("item-1", "2026-04-08", undefined);
    expect(cancelBookingMock).toHaveBeenCalledWith("existing-1");
    expect(wrapper.find('[data-cy="fp-booking-success"]').exists()).toBe(true);
  });

  it("ignores an on-behalf booking in the same area/day (self-scoped guard, story 36.9)", async () => {
    // for_user_id set: the current user made this for a colleague. for_user_name
    // may be missing when display-name lookup misses; it still must not be
    // treated as the user's own conflict.
    fetchMyBookingsMock.mockResolvedValue({
      data: [
        {
          id: "existing-1",
          type: "bookings",
          attributes: {
            item_id: "other-desk",
            item_name: "Desk Z",
            area_id: "area-1",
            booking_date: "2026-04-08",
            item_group_id: "ig-1",
            item_group_name: "Cube 1",
            area_name: "Office",
            created_at: "",
            booked_by_user_id: "me",
            booked_by_user_name: "Me",
            booked_for_me: false,
            for_user_id: "u-1",
            note: "",
          },
        },
      ],
    } as never);

    const wrapper = mountComponent({ areaId: "area-1" });
    await flushPromises();
    createBookingMock.mockClear();
    cancelBookingMock.mockClear();

    await wrapper.get('[data-cy="fp-item-item-1"]').trigger("click");
    await flushPromises();
    await wrapper.get('[data-cy="fp-booking-confirm"]').trigger("click");
    await flushPromises();

    // No swap dialog; books directly without cancelling the colleague's seat.
    expect(wrapper.find('[data-cy="confirm-dialog-confirm"]').exists()).toBe(false);
    expect(cancelBookingMock).not.toHaveBeenCalled();
    expect(createBookingMock).toHaveBeenCalledWith("item-1", "2026-04-08", undefined);
  });

  it("cancels only old floor-plan swap bookings whose multi-day replacement was created", async () => {
    fetchMyBookingsMock.mockResolvedValue({
      data: [
        {
          id: "existing-8",
          type: "bookings",
          attributes: {
            item_id: "other-desk-8",
            item_name: "Desk Y",
            area_id: "area-1",
            booking_date: "2026-04-08",
            item_group_id: "ig-1",
            item_group_name: "Cube 1",
            area_name: "Office",
            created_at: "",
            booked_by_user_id: "me",
            booked_by_user_name: "Me",
            booked_for_me: true,
            note: "",
          },
        },
        {
          id: "existing-9",
          type: "bookings",
          attributes: {
            item_id: "other-desk-9",
            item_name: "Desk Z",
            area_id: "area-1",
            booking_date: "2026-04-09",
            item_group_id: "ig-1",
            item_group_name: "Cube 1",
            area_name: "Office",
            created_at: "",
            booked_by_user_id: "me",
            booked_by_user_name: "Me",
            booked_for_me: true,
            note: "",
          },
        },
      ],
    } as never);
    createMultiDayBookingMock.mockResolvedValue({
      created: [
        {
          id: "new-8",
          type: "bookings",
          attributes: { item_id: "item-1", booking_date: "2026-04-08" },
        },
      ],
      conflicts: ["2026-04-09"],
    } as never);

    const wrapper = mountComponent({ areaId: "area-1" });
    await flushPromises();
    createMultiDayBookingMock.mockClear();
    cancelBookingMock.mockClear();

    await wrapper.get('[data-cy="fp-item-item-1"]').trigger("click");
    await flushPromises();
    await wrapper.get('[data-cy="fp-booking-day-TH"] input').trigger("change");
    await wrapper.get('[data-cy="fp-booking-confirm"]').trigger("click");
    await flushPromises();

    await wrapper.get('[data-cy="confirm-dialog-confirm"]').trigger("click");
    await flushPromises();
    await wrapper.get('[data-cy="confirm-dialog-confirm"]').trigger("click");
    await flushPromises();

    expect(createMultiDayBookingMock).toHaveBeenCalledWith(
      "item-1",
      ["2026-04-08", "2026-04-09"],
      undefined,
    );
    expect(cancelBookingMock).toHaveBeenCalledTimes(1);
    expect(cancelBookingMock).toHaveBeenCalledWith("existing-8");
    expect(cancelBookingMock).not.toHaveBeenCalledWith("existing-9");
  });

  it("does not prompt for a colleague booking when only the user has a conflict (story 36.9)", async () => {
    // The user has an own conflicting booking, but the NEW booking is for a
    // colleague — it never occupies the user's own seat, so no guard prompt.
    fetchColleaguesMock.mockResolvedValue({
      data: [{ id: "u-1", type: "colleagues", attributes: { display_name: "Jane Doe" } }],
    } as never);
    fetchMyBookingsMock.mockResolvedValue({
      data: [
        {
          id: "existing-1",
          type: "bookings",
          attributes: {
            item_id: "other-desk",
            item_name: "Desk Z",
            area_id: "area-1",
            booking_date: "2026-04-08",
            item_group_id: "ig-1",
            item_group_name: "Cube 1",
            area_name: "Office",
            created_at: "",
            note: "",
          },
        },
      ],
    } as never);

    const wrapper = mountComponent({ areaId: "area-1" });
    await flushPromises();
    createBookingMock.mockClear();
    cancelBookingMock.mockClear();
    fetchMyBookingsMock.mockClear();

    await wrapper.get('[data-cy="fp-item-item-1"]').trigger("click");
    await flushPromises();
    (wrapper.vm as unknown as { bookingColleagueId: string | null }).bookingColleagueId = "u-1";
    await flushPromises();
    await wrapper.get('[data-cy="fp-booking-confirm"]').trigger("click");
    await flushPromises();

    expect(wrapper.find('[data-cy="confirm-dialog-confirm"]').exists()).toBe(false);
    expect(cancelBookingMock).not.toHaveBeenCalled();
    expect(createBookingMock).toHaveBeenCalledWith("item-1", "2026-04-08", {
      forUserId: "u-1",
      forUserName: "Jane Doe",
    });
  });

  it("prompts when the colleague already has a booking in the same area/day (story 36.9)", async () => {
    fetchColleaguesMock.mockResolvedValue({
      data: [{ id: "u-1", type: "colleagues", attributes: { display_name: "Jane Doe" } }],
    } as never);
    // The user already booked Desk Z for colleague u-1 in the same area/day.
    fetchMyBookingsMock.mockResolvedValue({
      data: [
        {
          id: "existing-1",
          type: "bookings",
          attributes: {
            item_id: "other-desk",
            item_name: "Desk Z",
            area_id: "area-1",
            booking_date: "2026-04-08",
            item_group_id: "ig-1",
            item_group_name: "Cube 1",
            area_name: "Office",
            created_at: "",
            booked_by_user_id: "me",
            for_user_id: "u-1",
            for_user_name: "Jane Doe",
            note: "",
          },
        },
      ],
    } as never);

    const wrapper = mountComponent({ areaId: "area-1" });
    await flushPromises();
    createBookingMock.mockClear();
    cancelBookingMock.mockClear();

    await wrapper.get('[data-cy="fp-item-item-1"]').trigger("click");
    await flushPromises();
    (wrapper.vm as unknown as { bookingColleagueId: string | null }).bookingColleagueId = "u-1";
    await flushPromises();
    await wrapper.get('[data-cy="fp-booking-confirm"]').trigger("click");
    await flushPromises();

    // The colleague-variant swap prompt names the colleague; nothing created yet.
    const dialog = wrapper.find('[data-cy="confirm-dialog-confirm"]');
    expect(dialog.exists()).toBe(true);
    expect(wrapper.text()).toContain("Jane Doe");
    expect(createBookingMock).not.toHaveBeenCalled();

    // Confirming swaps: create the new booking first, then cancel the old one.
    await dialog.trigger("click");
    await flushPromises();
    expect(createBookingMock).toHaveBeenCalledWith("item-1", "2026-04-08", {
      forUserId: "u-1",
      forUserName: "Jane Doe",
    });
    expect(cancelBookingMock).toHaveBeenCalledWith("existing-1");
  });

  it("skips only the declined day in a multi-day swap (per-day guard, story 36.9)", async () => {
    // Two future days selected; an existing own booking conflicts on 2026-04-10
    // only. Declining that day must still book 2026-04-09.
    fetchMyBookingsMock.mockResolvedValue({
      data: [
        {
          id: "existing-1",
          type: "bookings",
          attributes: {
            item_id: "other-desk",
            item_name: "Desk Z",
            area_id: "area-1",
            booking_date: "2026-04-10",
            item_group_id: "ig-1",
            item_group_name: "Cube 1",
            area_name: "Office",
            created_at: "",
            booked_by_user_id: "me",
            booked_by_user_name: "Me",
            booked_for_me: true,
            note: "",
          },
        },
      ],
    } as never);

    const wrapper = mountComponent({ areaId: "area-1", selectedDay: "2026-04-09" });
    await flushPromises();
    createBookingMock.mockClear();
    createMultiDayBookingMock.mockClear();
    cancelBookingMock.mockClear();

    await wrapper.get('[data-cy="fp-item-item-1"]').trigger("click");
    await flushPromises();
    // Select 2026-04-10 in addition to the preselected 2026-04-09.
    await wrapper.get('[data-cy="fp-booking-day-FR"] input').trigger("change");

    await wrapper.get('[data-cy="fp-booking-confirm"]').trigger("click");
    await flushPromises();

    // Decline the swap for the conflicting day (2026-04-10).
    await wrapper.get('[data-cy="confirm-dialog-cancel"]').trigger("click");
    await flushPromises();

    // Only the non-conflicting day is booked; the declined day is skipped and
    // no cancellation happens.
    expect(createBookingMock).toHaveBeenCalledWith("item-1", "2026-04-09", undefined);
    expect(createMultiDayBookingMock).not.toHaveBeenCalled();
    expect(cancelBookingMock).not.toHaveBeenCalled();
  });

  it("leaves both bookings unchanged when the swap is cancelled (story 36.9)", async () => {
    fetchMyBookingsMock.mockResolvedValue({
      data: [
        {
          id: "existing-1",
          type: "bookings",
          attributes: {
            item_id: "other-desk",
            item_name: "Desk Z",
            area_id: "area-1",
            booking_date: "2026-04-08",
            item_group_id: "ig-1",
            item_group_name: "Cube 1",
            area_name: "Office",
            created_at: "",
            booked_by_user_id: "me",
            booked_by_user_name: "Me",
            booked_for_me: true,
            note: "",
          },
        },
      ],
    } as never);

    const wrapper = mountComponent({ areaId: "area-1" });
    await flushPromises();
    createBookingMock.mockClear();
    cancelBookingMock.mockClear();

    await wrapper.get('[data-cy="fp-item-item-1"]').trigger("click");
    await flushPromises();
    await wrapper.get('[data-cy="fp-booking-confirm"]').trigger("click");
    await flushPromises();

    await wrapper.get('[data-cy="confirm-dialog-cancel"]').trigger("click");
    await flushPromises();

    expect(cancelBookingMock).not.toHaveBeenCalled();
    expect(createBookingMock).not.toHaveBeenCalled();
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
