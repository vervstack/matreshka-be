import {useEffect, useState} from "react";
import {ListConfigsRequest, ListConfigsResponse} from "@vervstack/matreshka";

import cls from '@/pages/HomePage.module.css';

import ConfigInfoCard from "@/components/home/ConfigInfoCard.tsx";

import {useApi} from "@/app/hooks/api/api.ts";

import LoaderWrapper from "@/segments/LoaderWrapper.tsx";
import NewConfigCard from "@/components/home/NewConfigCard.tsx";

export default function HomePage() {
    const {ListConfigs} = useApi()

    const [list, setList] = useState<ListConfigsResponse | undefined>();

    const [loadFunc, setLoadFunc] = useState<Promise<void> | undefined>(undefined);

    useEffect(() => {
        const req = {} as ListConfigsRequest;
        setLoadFunc(
            ListConfigs(req)
                .then(setList)
        )
    }, []);

    return (
        <div className={cls.HomePageContainer}>
            <LoaderWrapper
                load={loadFunc}
            >
                <div className={cls.CardWrapper}>

                    {
                        list?.configs &&
                        list?.configs?.map((config, idx) => {
                            return (
                                <ConfigInfoCard
                                    key={idx}
                                    cardTitle={config.name || ''}/>
                            )
                        })
                    }

                    <NewConfigCard/>
                </div>


            </LoaderWrapper>
        </div>
    )
}
