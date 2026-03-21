import cn from "classnames";
import {useEffect, useState} from "react";
import {useSearchParams} from "react-router-dom";
import {
    Format,

    GetConfigRequest,
    GetConfigNodeRequest,
    SaveConfigRequest,
} from "@vervstack/matreshka";

import {useApi} from "@/app/hooks/api/api.ts";
import cls from "@/widget/ConfigDisplayWidget.module.css";
import Pencil from "@/assets/icons/Pencil.svg";
import X from "@/assets/icons/X.svg";
import Eye from "@/assets/icons/Eye.svg";

import LoaderWrapper from "@/segments/LoaderWrapper.tsx";
import SelectOneOf from "@/components/shared/SelectOneOf.tsx";
import SelectDropDown from "@/components/shared/SelectDropDown.tsx";


interface ConfigDisplayWidgetProps {
    configName: string;
}

export default function ConfigDisplayWidget({configName}: ConfigDisplayWidgetProps) {
    const {GetConfig, GetConfigNodes, SaveConfig} = useApi();

    const [searchParams, setSearchParams] = useSearchParams();
    const versionParam = searchParams.get("version");

    const [format, setFormat] = useState<Format>(Format.yaml);
    const [version, setVersion] = useState<string>(versionParam || "");
    const [versions, setVersions] = useState<string[]>([]);

    const [configText, setConfigText] = useState<string>("");
    const [editedConfigText, setEditedConfigText] = useState<string>("");


    const [isEditing, setIsEditing] = useState<boolean>(false);
    const [isComparing, setIsComparing] = useState<boolean>(false);
    const [configType, setConfigType] = useState<string>("");

    const originalLines = configText.split("\n");
    const editedLines = editedConfigText.split("\n");

    const lineCount = (isComparing ? Math.max(originalLines.length, editedLines.length) : (isEditing ? editedLines.length : originalLines.length));
    const editorHeight = `${lineCount * 1.5}em`;

    const [loadFunc, setLoadFunc] = useState<Promise<void> | undefined>(undefined);

    const loadConfig = () => {
        const req: GetConfigRequest = {
            configName: configName,
            format: format,
            version: version,
        };

        const promise = GetConfig(req)
            .then(res => {
                if (res.config) {
                    const text = atob(res.config.toString());
                    setConfigText(text);
                    setEditedConfigText(text);
                }

                res.info?.configBase?.configType && setConfigType(res.info.configBase.configType)
            });

        setLoadFunc(promise);
    };

    useEffect(() => {
        const req: GetConfigNodeRequest = {
            configName: configName,
        };
        GetConfig(req).then(res => {
            if (res.info?.versions) {
                setVersions([...res.info.versions]);

                if (!res.info.versions.includes(version)) {
                    setVersion(res.info.versions[0])
                }
            }
        }).catch(err => {
            console.error("Failed to fetch versions", err);
        });
    }, [configName, GetConfigNodes]);

    useEffect(() => {
        if (version !== "") {
            setSearchParams({version}, {replace: true});
        }
    }, [version, setSearchParams]);

    useEffect(() => {
        loadConfig();
    }, [configName, format, version, GetConfig]);

    useEffect(() => {
        document.title = configName;
        if (version != "") {
            document.title += `@${version}`;
        }
    }, [configName, version]);

    const isDirty = configText !== editedConfigText;

    const onSave = () => {
        const req: SaveConfigRequest = {
            configName: configName,
            format: format,
            version: version,
            config: new TextEncoder().encode(editedConfigText),
        };

        SaveConfig(req).then(() => {
            setConfigText(editedConfigText);
            setIsEditing(false);
        }).catch(err => {
            console.error("Failed to save config", err);
        });
    };

    const onRollback = () => {
        setEditedConfigText(configText);
        setIsEditing(false);
    };

    const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
        if (e.key === 'Tab') {
            e.preventDefault();
            const textarea = e.currentTarget;
            const start = textarea.selectionStart;
            const end = textarea.selectionEnd;

            const newValue = editedConfigText.substring(0, start) + "    " + editedConfigText.substring(end);
            setEditedConfigText(newValue);

            // Need to set cursor position after react re-render
            setTimeout(() => {
                textarea.selectionStart = textarea.selectionEnd = start + 4;
            }, 0);
        }
    };

    return (
        <div className={cls.ConfigDisplayWidgetContainer}>
            <div className={cls.Header}>
                <div className={cls.TitleSection}>
                    <div className={cls.ConfigName}>{configName}</div>
                    {configType && <div className={cls.ConfigType}>{configType}</div>}
                </div>
                <div className={cls.ActionsSections}>
                    <div className={cls.VersionSelect}>
                        <SelectDropDown
                            label="Version"
                            options={versions}
                            selectedOption={version}
                            onSelected={setVersion}
                        />
                    </div>
                    <div className={cls.FormatSelect}>
                        <SelectOneOf
                            options={[Format.env, Format.yaml]}
                            selectedOption={format}
                            onSelected={(opt) => setFormat(opt as Format)}
                        />
                    </div>
                </div>
            </div>

            <div className={cn(cls.EditSections,
                {
                    [cls.hidden]: isDirty,
                })}>
                <button className={cls.RollbackButton}
                        onClick={onRollback}>
                    rollback
                </button>
                <button className={cls.SaveButton}
                        onClick={onSave}>
                    save
                </button>
            </div>

            <LoaderWrapper load={loadFunc} onClickTryAgain={loadConfig}>
                <div className={cls.Content}>
                    <ContentActions
                        isEditing={isEditing}
                        isChanged={configText !== editedConfigText}
                        setIsComparing={setIsComparing}
                        isComparing={isComparing}
                        setIsEditing={setIsEditing}
                    />

                    <div className={cn(cls.ConfigBlocks, {[cls.Comparing]: isComparing})}>
                        {isComparing && (
                            <div className={cls.Column}>
                                <div className={cls.ColumnLabel}>Original</div>
                                <pre className={cls.Pre} style={{height: editorHeight}}>{configText}</pre>
                            </div>
                        )}
                        <div className={cls.Column}>
                            {isComparing && <div className={cls.ColumnLabel}>Current</div>}
                            {isEditing ? (
                                <textarea
                                    className={cn(cls.Editor, {[cls.Changed]: isComparing && configText !== editedConfigText})}
                                    value={editedConfigText}
                                    onChange={(e) => setEditedConfigText(e.target.value)}
                                    onKeyDown={handleKeyDown}
                                    style={{height: editorHeight}}
                                    autoFocus
                                />
                            ) : (
                                <pre
                                    className={cn(cls.Pre, {[cls.Changed]: isComparing && configText !== editedConfigText})}
                                    style={{height: editorHeight}}>{editedConfigText}</pre>
                            )}
                        </div>
                    </div>
                </div>
            </LoaderWrapper>
        </div>
    );
}

interface ContentActionsProps {
    isEditing: boolean;
    isChanged: boolean;

    setIsComparing: (v: boolean) => void;
    isComparing: boolean;

    setIsEditing: (v: boolean) => void;
}

function ContentActions({
                            isEditing, isChanged,
                            setIsComparing, isComparing,
                            setIsEditing
                        }: ContentActionsProps) {
    return (
        <div className={cls.TopActions}>
            {isEditing && isChanged && (
                <div className={cls.IconAction}
                     onClick={() => setIsComparing(!isComparing)}
                     data-tooltip-id="root-tooltip"
                     data-tooltip-content={isChanged ? "Hide original" : "Show original"}
                >
                    <img src={Eye} alt="Toggle Comparison"
                         style={{opacity: isComparing ? 1 : 0.5}}/>
                </div>)
            }

            <div className={cls.IconAction}
                 onClick={() => setIsEditing(!isEditing)}>
                {isEditing ?
                    (<img src={X} alt="Stop editing"/>)
                    :
                    (<img src={Pencil} alt="Edit"/>)}
            </div>
        </div>)
}
