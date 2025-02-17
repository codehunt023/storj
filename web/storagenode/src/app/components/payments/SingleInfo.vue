// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

<template>
    <div class="info-container" :style="{ width }">
        <div class="info-container__row">
            <p class="info-container__label">{{ label }}</p>
            <div class="info-container__info-area" v-if="infoText">
                <ChecksInfoIcon class="checks-area-image" alt="Info icon with question mark" @mouseenter="toggleTooltipVisibility" @mouseleave="toggleTooltipVisibility"/>
                <div class="tooltip" v-show="isTooltipVisible">
                    <div class="tooltip__text-area">
                        <p class="tooltip__text-area__text">{{ infoText }}</p>
                    </div>
                    <div class="tooltip__footer"></div>
                </div>
            </div>
        </div>
        <p class="info-container__value">{{ value }}</p>
    </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';

import ChecksInfoIcon from '@/../static/images/checksInfo.svg';

@Component ({
    components: {
        ChecksInfoIcon,
    },
})
export default class SingleInfo extends Vue {
    @Prop({default: '100%'})
    public readonly width: string;
    @Prop({default: 'Label'})
    public readonly label: string;
    @Prop({default: 'value'})
    public readonly value: string;
    @Prop({default: ''})
    private readonly infoText: string;

    /**
     * Indicates if tooltip needs to be shown.
     */
    public isTooltipVisible = false;

    /**
     * Toggles tooltip visibility.
     */
    public toggleTooltipVisibility(): void {
        this.isTooltipVisible = !this.isTooltipVisible;
    }
}
</script>

<style scoped lang="scss">
    .info-container {
        width: 45%;
        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: space-between;
        background: var(--block-background-color);
        border: 1px solid var(--block-border-color);
        box-sizing: border-box;
        border-radius: 10px;
        padding: 20px 30px;

        &__row {
            display: flex;
            align-items: center;
            justify-content: flex-start;
        }

        &__info-area {
            position: relative;
            width: 18px;
            height: 18px;
            margin-left: 10px;
        }

        &__label,
        &__value {
            font-family: 'font_regular', sans-serif;
            font-size: 14px;
            line-height: 20px;
            color: var(--regular-text-color);
        }

        &__value {
            font-family: 'font_bold', sans-serif;
            font-size: 20px;
        }
    }

    .tooltip {
        position: absolute;
        bottom: 35px;
        left: 50%;
        transform: translate(-50%);
        height: auto;
        box-shadow: 0 2px 48px var(--tooltip-shadow-color);
        border-radius: 12px;
        background: var(--tooltip-background-color);

        &__text-area {
            padding: 15px 11px;
            width: 360px;
            font-family: 'font_regular', sans-serif;
            font-size: 11px;
            line-height: 17px;
            color: var(--regular-text-color);
            text-align: center;
        }

        &__footer {
            position: absolute;
            left: 50%;
            transform: translate(-50%);
            width: 0;
            height: 0;
            border-style: solid;
            border-width: 11.5px 11.5px 0 11.5px;
            border-color: var(--tooltip-background-color) transparent transparent transparent;
        }
    }

    @media screen and (max-width: 460px) {

        .checks-area-image {
            display: none;
        }
    }
</style>
