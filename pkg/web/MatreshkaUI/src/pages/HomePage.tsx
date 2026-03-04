import {useEffect, useState} from "react";
import {ListConfigsRequest, ListConfigsResponse} from "@vervstack/matreshka";

import cls from '@/pages/HomePage.module.css';

import Card from "@/components/home/Card.tsx";

import {useApi} from "@/app/hooks/api/api.ts";

import LoaderWrapper from "@/segments/LoaderWrapper.tsx";

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
                {
                    list?.configs ? <div className={cls.CardWrapper}>
                        {
                            list?.configs?.map(config => {
                                return (
                                    <Card cardTitle={config.name || ''}/>
                                )
                            })
                        }
					</div> : <div>No configs</div>
                }
            </LoaderWrapper>
        </div>
    )
}
