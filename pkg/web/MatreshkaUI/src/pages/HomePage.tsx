import {useEffect, useState} from "react";
import {ListConfigsRequest, ListConfigsResponse} from "@vervstack/matreshka";

import cls from '@/pages/HomePage.module.css';

import Card from "@/components/home/Card.tsx";

import {useApi} from "@/app/hooks/api/api.ts";
import {useToaster} from "@/app/hooks/toaster/Toaster.ts";

export default function HomePage() {
    const {ListConfigs} = useApi()
    const toaster = useToaster();

    const [list, setList] = useState<ListConfigsResponse | undefined>();

    useEffect(() => {
        const req = {} as ListConfigsRequest;

        ListConfigs(req)
            .then(setList)
            .catch(toaster.catchGrpc)
    }, []);

    return (
        <div className={cls.HomePageContainer}>
            <div className={cls.CardWrapper}>
                {
                    list?.configs?.map(config => {
                        return (
                            <Card cardTitle={config.name || ''}/>
                        )
                    })
                }
            </div>
        </div>
    )
}
