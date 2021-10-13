///////////////////////////////////////////////////////////////////////
//                                                                   // 
//    Copyright� 2014 Elo Touch Solutions. All rights reserved.      //
//                                                                   //
//    This computer program source code is provided "as is"          //
//    without warranty of fitness or suitability for any purpose.    //
//    Elo Touch Solutions is not responsible for consequential or    //
//    incidental damages resulting from or related to the use of     //
//    this code. This code may be used in part or in its entirety    //
//    only to support the use of Elo Touch Solutions products.       //
//                                                                   //
///////////////////////////////////////////////////////////////////////

// EloStructs.h

#include <wchar.h>
#include <stddef.h>
#include <stdbool.h>
#include <stdio.h>
#include <string.h>
#include <windows.h>

#ifndef ELO_STRUCTS_H
#define ELO_STRUCTS_H

#define MAX_ELO_MONITORS		16
#define MAX_ELO_TOUCHSCREENS	16 // MAXIMUM_WAIT_OBJECTS is 64

#define ELO_MT_MAX_COUNT 64 // maximum number of contacts (used by both apps and usb driver).

#define COMPORT_NAME_LENGTH					256
#define STACKDEVICE_NAME_LENGTH				256
#define SERIALNUMBER_NAME_LENGTH			32
#define USB_PRODUCT_STRING_LENGTH			64
#define SCF_HEADER_SERIAL_NUMBER_LENGTH		12
#define DATA_FILE_MAX_PATH					256
#define DATA_FILE_MAX_NAME					32
#define PROGRAM_FILE_MAX_NAME				64
#define USB_SERIAL_NUMBER_LENGTH			10

#define MAX_BOUND			64
#define MAX_BUTTON_BOUND    64

#define APR_SET_COUNT 4
#define APR_PARAM_COUNT 48

// Beep properties
#define DEF_BEEP_FREQUENCY 800 // Hz
#define DEF_BEEP_DURATION  100 // ms

// Untouch timeout
#define MIN_UNTOUCH_TIMEOUT 1   // in sec (1 sec)
#define MAX_UNTOUCH_TIMEOUT 300 // in sec (5 min)

// Drag delay
#define MIN_DRAG_DELAY_TIME 100  // in milli-sec (0.1 sec)
#define MAX_DRAG_DELAY_TIME 2000 // in milli-sec (2 sec)

// Anchor shift: Used by both Untouch timeout AND Drag Delay
#define DEF_ANCHOR_SHIFT 1 // mm 
#define MIN_ANCHOR_SHIFT 0 // mm
#define MAX_ANCHOR_SHIFT 10 // mm

#define MIN_BEEP_FREQUENCY 500 
#define MAX_BEEP_FREQUENCY 4000

#define MIN_BEEP_DURATION 20 
#define MAX_BEEP_DURATION 500 

typedef struct
{
	UCHAR        T;                        // Leading 'T' for touch packet
	UCHAR        Status;                    // Touch status
	SHORT        X;                        // X value
	SHORT        Y;                        // Y value
	SHORT        Z;                        // Z value
} TOUCHPACKET, *PTOUCHPACKET;

typedef enum tagTouchStatus   // the blocked call returns data depending on this flag
{
	InitialTouch = 1,
	StreamTouch = 2,
	UnTouch = 4
} TOUCH_STATUS;

typedef	struct tagTouch
{
	TOUCH_STATUS status;
	UCHAR  id;     // touch id: 0 or 1 representing the touch finger
	USHORT x;      // x of a touch point
	USHORT y;      // y of a touch point      
} TOUCH, *PTOUCH;

typedef struct
{
	TOUCH touch[ELO_MT_MAX_COUNT]; // max number of touches (fingers)
	int count; // actual touch count 
} MT_TOUCH, *PMT_TOUCH;

typedef enum _GETPOINTS_CODE    // the blocked call returns data depending on this flag
{
	NoRequest = 0,
	ReturnImmediately,
	ReturnOnTouch,
	ReturnOnUntouch,
	ReturnOnNextValidTouch,

}GETPOINTS_CODE;

typedef enum _GETPOINTS_TRANSLATION
{
	Raw = 1,
	Translated,
	ControllerDirect
}GETPOINTS_TRANSLATION;

typedef struct _ELO_GETTOUCHPOINTS
{
	BOOL                    ValidFlag;    // OUT-if TRUE, points are valid untouch points 
	GETPOINTS_TRANSLATION    Translation;// OUT-if TRUE, points are translated
	LONG                    XValue;        // OUT- x value of touch
	LONG                    YValue;        // OUT- x value of touch
	LONG                    ZValue;        // OUT- Z value of touch
	TOUCHPACKET                Controller_Packet;
	TOUCH_STATUS            Status;
	GETPOINTS_CODE    GetPointsFlag; // IN - the blocked call returns data depending on this flag
} ELO_GET_TOUCH_POINTS;

typedef struct _AccelBounds		// IN Set calibrate data input buffer
{
	float	X_Max;
	float	X_Min;
	float	Y_Max;
	float	Y_Min;
	float	Z_Max;
	float	Z_Min;
} AccelBounds, *pAccelBonds;

typedef struct _ACCELDATA
{
	ULONG		Enable;
	ULONG		Scale;
	AccelBounds	Bounds;
} ELO_ACCEL_DATA;

#ifndef _CONTRL_STAT_DEF
#define _CONTRL_STAT_DEF
typedef enum _CONTRL_STAT                // ctrl_status values
{
	CS_OK = 0,
	CS_ConstantTouch,
	CS_CanNotFindController,
	CS_NoResponse,
	CS_InvalidResponse,
	CS_CanNotSetBaudRate,
	CS_CommandNotSent,
	CS_SystemError,
	CS_InvalidCommPort,
	CS_CommPortFailedOpen,
	CS_CommPortCommandError,
	CS_CommPortNoController,
	CS_UndefinedController
} CONTRL_STAT;
#endif

typedef enum
{
	enumVrtlDeskDisabled = 0,    // no virtual desktop
	enumVrtlNoBounds,            // NO virtual desktop bounds, No Clipping, Cursor visible
	enumVrtlBoundsClipped,       // virtual desktop bounds enable, Clipping, Cursor moves at bounds
	enumVrtlBoundsFreeze        // virtual desktop bounds enable, Clipping, Cursor Freezes at bounds
} VrtlBoundMode;

typedef struct _ClippingBounds        // IN Set calibrate data input buffer
{
	long            X_Min;
	long            X_Max;
	long            Y_Min;
	long            Y_Max;
	long            Z_Min;
	long            Z_Max;
} ClippingBounds, *PClippingBounds;


typedef struct _EloCalData        // IN Set calibrate data input buffer
{
	LONG            EloDx;
	LONG            ScrDx;
	LONG            X_Offset;
	LONG            EloDy;
	LONG            ScrDy;
	LONG            Y_Offset;
	LONG            EloDz;
	LONG            ScrDz;
	LONG            Z_Offset;
	LONG            xyswap;

	ULONG           EloMonitorNumber;
	ULONG           MonitorSerialNumber;
	unsigned short	Checksum;
	int				nScreenIndex;
	unsigned char   controllerMode[2];

	// Support multiple screens for mouse collection
	LONG            xVirtScrSize;
	LONG            yVirtScrSize;
	LONG            xVirtScrCoord;
	LONG            yVirtScrCoord;

} ELO_CAL_DATA;

typedef struct _DIAGNOSTICS
{
	CONTRL_STAT      ctrl_status;         // OUT- Controller Status 
	ULONG            HardwareHandShaking; // OUT-Hardware handshaking turned on /off
	LONG             BaudRate;            // OUT- Baud rate of controler, 0 for bus

	unsigned char    crevmajor;           // OUT- controller rev major number
	unsigned char    crevminor;           // OUT- controller rev minor number
	unsigned char    crevbuild;

	unsigned char    trevmajor;           // OUT- Unused
	unsigned char    trevminor;           // OUT- Unused
	unsigned char    trevbuild;

	unsigned char    diagcodes[8];        // OUT- Diag codes ret from controller
	unsigned char    id[8];               // OUT- OEM ID string ret from controller
	unsigned char    cnt_id[8];           // OUT- Full Smartset controller ID packet 
	unsigned char    driver_id[32];       // OUT- Driver ID

	unsigned char    hidReportId;

	// add on fields from enum ioctl
	ULONG   uInterfaceType;                           // OUT- TOUCHSCREEN_TYPE_USB, TOUCHSCREEN_NT_SERIAL, ...
	wchar_t PortFriendlyName[COMPORT_NAME_LENGTH];   // OUT- Used for serial touchscreens
	wchar_t SerialNumber[SERIALNUMBER_NAME_LENGTH];  // OUT- ASCII 8-digit serial#
	wchar_t wcUsbProductString[USB_PRODUCT_STRING_LENGTH];

	// APR only
	char szCalFileName[DATA_FILE_MAX_NAME];
	int nCalFileVerMajor;
	int nCalFileVerMinor;
	char szSensorSN[SCF_HEADER_SERIAL_NUMBER_LENGTH];
	// End APR only
} ELO_DIAGNOSTICS;

typedef struct _ELO_CLIPPING_MODE
{
	VrtlBoundMode    ClippingMode;
	ULONG            NumBounds;
	ULONG            ExclusionFlag;
	ClippingBounds    Bounds[MAX_BOUND];
	int                MonitorNumber;
} ELO_CLIPPING_MODE;

typedef struct _ELO_FULL_SCREEN_CLIPPING
{
	VrtlBoundMode    ClippingMode;
	ClippingBounds    Bounds[MAX_BOUND];
} ELO_FULL_SCREEN_CLIPPING;

typedef struct _ELO_BUTTON_SEQ
{
	USHORT            InitialTouchSeq[MAX_BUTTON_BOUND];
	ULONG            NumInitialTouchSeq;
	USHORT            StreamTouchSeq[MAX_BUTTON_BOUND];
	ULONG            NumStreamTouchSeq;
	USHORT            UnTouchSeq[MAX_BUTTON_BOUND];
	ULONG            NumUnTouchSeq;
} ELO_BUTTON_SEQ;

typedef struct _ELO_RIGHT_CLICK_ON_HOLD
{
	ULONG            RightClickHW;
	ULONG            InitialTimeout;
	ULONG            DefaultRightClickDelay;
	ULONG            MaxRightClickDelay;
	ULONG            MinRightClickDelay;
	ULONG            ClickCount;        // IN - Touch Count for enabling right click on hold feature 
	ULONG            Active;        // IN /OUT Saves if this feature should be turned on , on reboot 
} ELO_RIGHT_CLICK_ON_HOLD;

// -----------------------------------------------------------------------------
//    MM_LOAD_USER_FILE
//

typedef enum _USER_FILE        // Direction of data in the ioctl call     
{
	SCF = 1,
	LUT,
	AUDIO
} USER_FILE;


typedef struct _ELO_LOAD_USER_FILE
{
	ULONG            Length;
	USER_FILE        UserFile;
	char            Filename[64];
	ULONG            CheckSum;
	ULONG            VerifyOnly;
	unsigned int    Buffer[1];
} ELO_LOAD_USER_FILE;

typedef enum _LIVE_SOUND_MODE
{
	LiveSoundSignal = 1,
	LiveSoundSpectrum,
	LiveSoundMatchSpectrum,
	LiveSoundBalanceData
}LIVE_SOUND_MODE;

typedef struct _ELO_LIVE_SOUND
{
	LIVE_SOUND_MODE    Mode;
	LONG            Length;
	LONG            Channels;
	LONG            DriverResidual;
	char            *DriverPointer;
	UCHAR            Buffer[4];
} ELO_LIVE_SOUND;

typedef enum __MODE_SWITCH
{
	eUseHashTests = 1,
	eUseLocalSearch
}MODE_SWITCH;

typedef enum _AUDIO_RECORED	// Direction of data in the ioctl call     
{
	RECORD_SIG = 1,
	RECORD_LVS,
	RECORD_ALL
} AUDIO_RECORD;

typedef enum MONITOR_ORIENTATION_TAG
{	// These are the same as defined by Windows
	moLandscape = 0,		// DMDO_DEFAULT
	moPortrait = 1,			// DMDO_90, rotate 90 degrees counter-clockwise
	moLandscapeFlipped = 2,	// DMDO_180, rotate 180 defrees 
	moPortraitFlipped = 3,	// DMDO_270, rotate 270 degrees counter-clockwise
} MONITOR_ORIENTATION;

typedef struct SPAN_MODE_TAG
{
	int row_index;
	int col_index;
	int tot_rows;
	int tot_cols;
	MONITOR_ORIENTATION orientation;
} SPAN_MODE;

typedef struct _ELO_GET_SERIAL_NUMBERS
{
	WCHAR  UsbSerialNumber[USB_SERIAL_NUMBER_LENGTH];
	CHAR   SensorSerialNumber[SCF_HEADER_SERIAL_NUMBER_LENGTH];
	CHAR   CalFileInUse[DATA_FILE_MAX_NAME];
	UCHAR  Diagnostics;
	LONG   CalFileMajorVerion;
	LONG   CalFileMinorVerion;
} ELO_GET_SERIAL_NUMBERS;

typedef struct _AUTO_SIZE
{
	// ToDriver
	BOOL          bDoAutoSizing;
	BOOL          bLockAutoSizing; // if bDoAutoSizing is true, at end of auto
								   // sizing status will be set to bLockAutoSizing
	// FromDriver
	BOOL          bIsLocked;
	USHORT        CycleCount;
	USHORT        DurationX;
	USHORT        DurationY;
} ELO_AUTO_SIZE;

typedef enum ELO_TOUCH_SENSITIVITY_TAG // the blocked call returns data depending on this flag
{
	Enhanced = 0,
	Normal = 1,
} ELO_TOUCH_SENSITIVITY;

typedef enum tagBeepSource
{
	BeepOff = 0x00, // 0000
	BeepBeeper = 0x01, // 0001
	BeepSpeaker = 0x02, // 0010
	BeepIRMonitor = 0x04, // 0100
} BEEP_SOURCE;

typedef enum tagCleanScreenMethod
{
	eNothing,
	eTimeElapsed,
	eCornerTouch,
	eGesture,
} ENUM_CLEAN_SCREEN_METHOD;

typedef enum tagCleanScreenTouchOptions
{
	eKeepTouch,
	eDisableTouch,
	eDisableController,
} ENUM_TOUCH_OPTION;

typedef struct ELO_BEEP_TAG
{
	ULONG       beepSource; // bit-wise or of BEEP_SOURCE fields
	int			nBeeperFrequency;
	int			nBeeperDuration;
	wchar_t		szSpeakerWaveFile[DATA_FILE_MAX_PATH];
} ELO_BEEP;

typedef struct ELO_RIGHT_BUTTON_TAG
{
	ULONG SwapCount; // IN - Touch Count for keeping button swapped
	ULONG Active;   // IN /OUT Saves if right button should be run on reboot 
} ELO_RIGHT_BUTTON;

typedef struct ELO_IR_BEAM_MONITOR_TAG
{
	BOOL	bEnable;
	BOOL	bLogToFile;
	int		nScanFrequency; // Time in sec between IR Beam test loop
} ELO_IR_BEAM_MONITOR;

typedef struct ELO_TOUCH_RESTRAINT_TAG
{
	LONG	nTimeLimit; // in seconds for untouch timeout, ms in drag delay
	LONG	nDistLimit; // in Elo controller Coord (0-4k).  
} ELO_TOUCH_RESTRAINT;

typedef struct ELO_MONITOR_RES_DRAG_THRESHOLD_TAG
{
	LONG xRes;		// monitor resolution in pixels
	LONG yRes;

	LONG xDragPix;  // System param for a drag
	LONG yDragPix;
} ELO_MONITOR_RES_DRAG_THRESHOLD;


typedef struct ELO_CLEAN_SCREEN_CONFIG_TAG
{
	ENUM_CLEAN_SCREEN_METHOD Method; // invocation methods
	ENUM_TOUCH_OPTION TouchOption; // touch options during cleaning session

	LONG nTimeElapsed; // For Time Elapsed invocation, length of time in minutes 
	LONG nTimeout; // cleaning session timeout if touch disabled (driver or ctrl disabled)
	LONG nIntervalStart; // minute count at start of interval if "Time Elapsed" used
	BOOLEAN bPause; // kernel action pasued if TRUE
	wchar_t szCleaningInst[DATA_FILE_MAX_PATH]; // user defined cleaning instructions

} ELO_CLEAN_SCREEN_CONFIG;

//#ifdef _WINDOWS

typedef struct _SmartsetData
{
	LONG ScreenNumber;

	union
	{
		char Command[8];
		unsigned short Response;
	};
} SMARTSETDATA, *PSMARTSETDATA;

typedef struct _SmartsetPkt
{
	unsigned char pkt[8];
} SMARTSET_PKT;

typedef struct MONITOR_TAG
{
	int		elo_mon_num;   // Elo's monitor number
	int		x;
	int		y;
	int		width;
	int		height;
	DWORD   orientation;
	DWORD   edid_sn;       // Monitor serial number from EDID 
	char    edid_mfr[4]; // null terminated 3-char manufacture ID
	int     edid_mon_size; // Monitor diagnal size (mm)
	unsigned char edid[1024];
	HMONITOR hMon;

	wchar_t device_id[MAX_PATH];
	BOOL    is_primary;
	SPAN_MODE span_mode;
} MONITOR;

typedef struct SCREEN_TAG
{
	int					nScreenIndex;		// device enumeration index
	USHORT				uVendorID;
	USHORT				uProductID;
	USHORT				uVersionNumber;
	wchar_t				szDevicePath[MAX_PATH];

	MONITOR*			pMonitor;			// pointer to the calibrated monitor

	ELO_ACCEL_DATA		accelData;
	ELO_BEEP			beepOptions;
	ELO_CAL_DATA		calData;
	ELO_CAL_DATA		onboardCalData;		// controller cal data
	ELO_CLIPPING_MODE	clipMode;
	ELO_DIAGNOSTICS		diag;
	ELO_IR_BEAM_MONITOR	beamMonitor;		// IR beam monitoring options
	ELO_RIGHT_BUTTON	rightBtn;
	ELO_CLEAN_SCREEN_CONFIG cleanScreen;

	LONG				lSmartSetToken;		// SmartSet token returned by the driver when handling SmartSet commands
	DWORD				dwSSTokenThreadId;	// ID of the thread which is holding SmartSet token

	HANDLE				hRBThread;
	HANDLE				hCalTouchThread;
	PVOID				pIrBeamMonitorDlg; // point to the IR Beam Monitoring Cross bar UI (CBeamHandler).
} SCREEN;

//#endif // _WINDOWS

#endif // ELO_STRUCTS_H