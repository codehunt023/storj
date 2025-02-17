// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

<template>
    <div class="credit-history__wrapper">
        <VButton
            class="credit-history__add-button"
            height="44px"
            width="174px"
            :on-press="onCreateClick"
            label="Add Coupon Code"
            v-if="couponCodeBillingUIEnabled"
        />
        <div class="credit-history__container">
            <div class="credit-history__content">
                <div class="credit-history__title-area">
                    <h1 class="credit-history__title">Credit History</h1>
                </div>
                <VLoader v-if="isHistoryFetching"/>
                <template v-else>
                    <SortingHeader/>
                    <CreditsItem
                        v-for="(item, index) in historyItems"
                        :key="index"
                        :credits-item="item"
                    />
                </template>
            </div>
        </div>
    </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';

import CreditsItem from '@/components/account/billing/freeCredits/CreditsItem.vue';
import SortingHeader from '@/components/account/billing/freeCredits/SortingHeader.vue';
import VButton from '@/components/common/VButton.vue';
import VLoader from '@/components/common/VLoader.vue';

import { RouteConfig } from '@/router';
import { PAYMENTS_ACTIONS } from '@/store/modules/payments';
import { PaymentsHistoryItem, PaymentsHistoryItemType } from '@/types/payments';

@Component({
    components: {
        CreditsItem,
        SortingHeader,
        VButton,
        VLoader,
    },
})
export default class CreditsHistory extends Vue {
    public isHistoryFetching = true;

    /**
     * Lifecycle hook after initial render.
     * Fetches payments history.
     */
    public async mounted(): Promise<void> {
        try {
            await this.$store.dispatch(PAYMENTS_ACTIONS.GET_PAYMENTS_HISTORY);

            this.isHistoryFetching = false;
        } catch (error) {
            await this.$notify.error(error.message);
        }
    }

    /**
     * Returns list of free credit history items.
     */
    public get historyItems(): PaymentsHistoryItem[] {
        return this.$store.state.paymentsModule.paymentsHistory.filter((item: PaymentsHistoryItem) => {
            return item.type === PaymentsHistoryItemType.Coupon;
        });
    }

    /**
     * Returns remaining sum of items.
     */
    public get remainingSum(): number {
        const remainingAmounts: number[] = this.historyItems.map((item: PaymentsHistoryItem) => item.remaining);

        return remainingAmounts.reduce((accumulator, current) => accumulator + current);
    }

    /**
    * Opens add coupon modal
    */
    public onCreateClick(): void {
        this.$router.push(RouteConfig.Billing.with(RouteConfig.AddCouponCode).path);
    }

    /**
     * Indicates if coupon code ui is enabled on the billing page.
     */
    public get couponCodeBillingUIEnabled(): boolean {
        return this.$store.state.appStateModule.couponCodeBillingUIEnabled;
    }
}
</script>

<style scoped lang="scss">
    p,
    h1 {
        margin: 0;
    }

    .credit-history {
        margin-top: 27px;
        padding: 0 0 80px 0;
        background-color: #f5f6fa;
        font-family: 'font_regular', sans-serif;

        &__wrapper {
            margin-bottom: 30px;
        }

        &__back-area {
            display: flex;
            align-items: center;
            cursor: pointer;
            width: 184px;
            margin-bottom: 32px;

            &__title {
                font-family: 'font_medium', sans-serif;
                font-weight: 500;
                font-size: 16px;
                line-height: 21px;
                color: #768394;
                white-space: nowrap;
                margin-left: 15px;
            }

            &:hover {

                .credit-history__back-area__title {
                    color: #2683ff;
                }

                .back-button-svg-path {
                    fill: #2683ff;
                }
            }
        }

        &__add-button {
            float: right;
            margin-bottom: 20px;
        }

        &__title {
            font-family: 'font_bold', sans-serif;
            font-size: 22px;
            line-height: 27px;
            color: #384b65;
            margin-bottom: 20px;
        }

        &__content {
            background-color: #fff;
            padding: 40px 40px 30px 40px;
            border-radius: 8px;
            display: flex;
            flex-direction: column;
            align-items: flex-start;
            clear: right;

            &__sum {
                font-family: 'font_bold', sans-serif;
                font-size: 36px;
                line-height: 53px;
                color: #384b65;
            }

            &__info {
                font-size: 16px;
                line-height: 24px;
                color: #909090;
                margin-bottom: 35px;
            }

            &__details {
                width: 100%;
                text-align: left;
                font-weight: 500;
                font-size: 16px;
                line-height: 23px;
                letter-spacing: 0.04em;
                color: #919191;
                padding-bottom: 22px;
                border-bottom: 1px solid #c7cdd2;
                margin-bottom: 75px;
            }
        }
    }

    ::-webkit-scrollbar,
    ::-webkit-scrollbar-track,
    ::-webkit-scrollbar-thumb {
        width: 0;
    }

    @media (max-height: 1000px) and (max-width: 1230px) {

        .credit-history {
            overflow-y: scroll;
            height: 65vh;
        }
    }
</style>
