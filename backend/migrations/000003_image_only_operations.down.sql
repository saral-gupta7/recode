ALTER TABLE jobs DROP CONSTRAINT jobs_operation_valid;

ALTER TABLE jobs ADD CONSTRAINT jobs_operation_valid CHECK (
    operation IN (
        'video_grayscale',
        'video_extract_audio',
        'video_remove_audio',
        'video_convert',
        'video_clip',
        'image_grayscale',
        'image_convert',
        'image_compress',
        'image_resize'
    )
) NOT VALID;
