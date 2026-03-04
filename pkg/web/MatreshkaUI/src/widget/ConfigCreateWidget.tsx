import {useState} from "react";
import {ConfigTypePrefix, CreateConfigRequest} from "@vervstack/matreshka";

import cls from "@/widget/ConfigCreateWidget.module.css";
import ConfigFile from "@/assets/icons/ConfigFile.svg";

import Input from "@/components/shared/Input.tsx";
import ActionButton from "@/components/shared/ActionButton.tsx";

import SelectDropDown from "@/components/shared/SelectDropDown.tsx";
import {useApi} from "@/app/hooks/api/api.ts";
import LoaderWrapper from "@/segments/LoaderWrapper.tsx";


export default function ConfigCreateWidget() {
    const [name, setName] = useState('');
    const [type, setType] = useState<ConfigTypePrefix>(ConfigTypePrefix.verv);

    const [func, setFunc] = useState<Promise<void> | undefined>();

    const {CreateConfig} = useApi();

    function create() {
        const req = {
            configName: name,
            configType: type,
        } as CreateConfigRequest;

        setFunc(CreateConfig(req)
            .then())
    }

    function clear() {
        setName("")
        setFunc(undefined)
    }

    return (
        <div className={cls.ConfigCreateWidgetContainer}>
            <LoaderWrapper
                load={func}
                onClickTryAgain={clear}
            >
                <div className={cls.Content}>
                    <div className={cls.Header}>Create new config</div>
                    <Input
                        label={'Name'}
                        inputValue={name}
                        onChange={setName}
                    />

                    <SelectDropDown
                        label={'Type'}
                        options={[
                            ConfigTypePrefix.plain, ConfigTypePrefix.verv, ConfigTypePrefix.minio,
                            ConfigTypePrefix.pg, ConfigTypePrefix.nginx, ConfigTypePrefix.kv,
                        ]}
                        selectedOption={type}
                        onSelected={(s: string) => setType(s as ConfigTypePrefix)}
                    />

                    <div className={cls.ButtonsWrapper}>
                        <ActionButton
                            tooltip={'Create new config'}
                            iconPath={ConfigFile}
                            alt={'Create'}
                            onClick={create}
                        />
                    </div>
                </div>
            </LoaderWrapper>
        </div>
    )
}
