import { File, Loader2, Paperclip } from 'lucide-react';
import { FileInput, FileUploader, FileUploaderContent, FileUploaderItem } from '../ui/file-upload';
import { useState } from 'react';
import OrganizationsApiService from '@/api/OrganizationsApiService';
import { toast } from '../ui/use-toast';

type Props = {
    organizationId: number;
};

const OrganizationDataUpload = ({ organizationId }: Props) => {
    const [files, setFiles] = useState<File[] | null>([]);
    const [isFileUploading, setIsFileUploading] = useState(false);

    const uploadFiles = (files: File[] | null) => {
        setFiles(files);

        const file = files?.[0];

        if (file) {
            setIsFileUploading(true);

            OrganizationsApiService.uploadFile(file, organizationId)
                .then(() => {
                    toast({
                        title: 'Успех',
                        description: 'Файл успешно загружен',
                        variant: 'default',
                    });
                })
                .catch(() => {
                    toast({
                        title: 'Успех',
                        description: 'Файл успешно загружен',
                        variant: 'default',
                    });
                })
                .finally(() => {
                    setIsFileUploading(false);
                });
        }
    };

    const dropzone = {
        accept: {
            'application/zip': ['.zip'],
        },
        multiple: false,
        maxFiles: 1,
    };

    return (
        <div className='mt-7'>
            <div className='flex flex-col'>
                <h1 className='font-semibold text-lg md:text-2xl'>Загрузка данных</h1>

                <p>Загрузите данные в формате .zip. Архив не должен содержать папки.</p>
            </div>

            <FileUploader
                value={files}
                onValueChange={uploadFiles}
                dropzoneOptions={dropzone}
                className='relative bg-background rounded-lg p-2 max-w-md'
            >
                <FileInput className='outline-dashed outline-1'>
                    <div className='flex items-center justify-center flex-col pt-3 pb-4 w-full '>
                        {isFileUploading ? (
                            <>
                                <Loader2 className='h-4 w-4 animate-spin' />
                            </>
                        ) : (
                            <File />
                        )}
                    </div>
                </FileInput>
                <FileUploaderContent>
                    {files &&
                        files.length > 0 &&
                        files.map((file, i) => (
                            <FileUploaderItem key={i} index={i}>
                                <Paperclip className='h-4 w-4 stroke-current' />
                                <span>{file.name}</span>
                            </FileUploaderItem>
                        ))}
                </FileUploaderContent>
            </FileUploader>
        </div>
    );
};

export default OrganizationDataUpload;
