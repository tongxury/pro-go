import cv2
import sys
import os

def get_frame(video_path, timestamp_ms, output_path):
    try:
        cap = cv2.VideoCapture(video_path)
        if not cap.isOpened():
            print(f"Error: Could not open video {video_path}", file=sys.stderr)
            return False

        # OpenCV 使用毫秒 (ms)
        cap.set(cv2.CAP_PROP_POS_MSEC, timestamp_ms)
        
        success, frame = cap.read()
        if success:
            # 写入临时文件，质量设为 95
            cv2.imwrite(output_path, frame, [int(cv2.IMWRITE_JPEG_QUALITY), 95])
            cap.release()
            return True
        else:
            # 如果指定位置失败，尝试取第一帧兜底
            cap.set(cv2.CAP_PROP_POS_MSEC, 0)
            success, frame = cap.read()
            if success:
                cv2.imwrite(output_path, frame, [int(cv2.IMWRITE_JPEG_QUALITY), 95])
                cap.release()
                return True
            
        cap.release()
        return False
    except Exception as e:
        print(f"Error: {str(e)}", file=sys.stderr)
        return False

if __name__ == "__main__":
    if len(sys.argv) < 4:
        sys.exit(1)
    
    video_path = sys.argv[1]
    # 输入是秒，转为毫秒
    timestamp_ms = float(sys.argv[2]) * 1000
    output_path = sys.argv[3]
    
    if get_frame(video_path, timestamp_ms, output_path):
        sys.exit(0)
    else:
        sys.exit(1)
