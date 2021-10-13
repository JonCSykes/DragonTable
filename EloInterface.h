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

// EloInterface.h
//
// The following ifdef block is the standard way of creating macros which make exporting 
// from a DLL simpler. All files within this DLL are compiled with the ELOMTAPI_EXPORTS
// symbol defined on the command line. this symbol should not be defined on any project
// that uses this DLL. This way any other project whose source files include this file see 
// ELOMT_API functions as being imported from a DLL, whereas this DLL sees symbols
// defined with this macro as being exported.

#pragma once

#include "EloStructs.h"

// The format of the version lock is two digits for each of the 3 version numbers (major/minor/build):
// For example, for release of v6.9.12, it must be in the format of 60912.

#define ELO_SDK_VERSION_LOCK 60912


// -----------------------------------------------------------------------------
// Interface function calls 
// -----------------------------------------------------------------------------

#ifdef ELOMTAPI_EXPORTS
#define ELOMT_API __declspec(dllexport)
#else
#define ELOMT_API __declspec(dllimport)
#endif

#ifdef __cplusplus
extern "C" {
#endif 

// Starting with v6.9.12, check whether caller's Version Locak (defined in EloInterface.h as
// ELO_SDK_VERSION_LOCK) is the same as the one from currently installed Elo multi-touch pacakge 
ELOMT_API int EloValidateSDKVersionLock (int nCallerVersionLock);

// Get total number of touch screens for current system
ELOMT_API int EloGetScreenCount();

// Enable or disable touch of specified screen by screen index
ELOMT_API bool EloEnableTouch (int nScreenIndex, bool bEnableTouch);

// Get touch state of specified by screen index
ELOMT_API bool EloIsTouchEnabled (int nScreenIndex);

// Get touch point from touch screen. *touch_status return touch status
ELOMT_API bool EloGetTouchPacket (int nScreenIndex, int* x, int* y, int* z, TOUCH_STATUS* touch_status, bool bRaw);

// Get calibration data for a touch screen
ELOMT_API bool EloGetCalData (int nScreenIndex, ELO_CAL_DATA * pCalData);

// Set calibration data for a touch screen (will alter controller calibration parameters)
ELOMT_API bool EloSetCalData (int nScreenIndex, ELO_CAL_DATA * pCalData);

// Set calibration data for a touch screen
ELOMT_API bool EloSetCalDataEx (int nScreenIndex, ELO_CAL_DATA * pCalData, bool bManualCalibration);

// Get diagnostics data for a specified touch screen
ELOMT_API DWORD EloGetDiagnosticsData(int nScreenIndex, ELO_DIAGNOSTICS * pDiag);

// Get the rectangle of current clipping area
ELOMT_API int EloGetClipRectangles(int nScreenIndex, ELO_CLIPPING_MODE* pClipMode);

// Set clipping rectangle bounds
ELOMT_API int EloSetClipRectangles(int nScreenIndex, ELO_CLIPPING_MODE* pClipMode);

// Get Edge Acceleration data for a touch screen
ELOMT_API int EloGetEdgeAcceleration(int nScreenIndex, ELO_ACCEL_DATA* pAccel);

// Set Edge Acceleration for a touch screen
ELOMT_API int EloSetEdgeAcceleration(int nScreenIndex, const ELO_ACCEL_DATA* pAccel);

// Get SCREEN struct by a specified screen by screen index
ELOMT_API SCREEN* EloGetScreenByIndex (int nScreenIndex);

// Flush the controller buffer by a specified screen by a screen index
ELOMT_API bool EloFlushControllerSmartsetBuffer (int nScreenIndex);

// Get controller firmware version
ELOMT_API bool EloGetControllerFWVersion (int nScreenIndex, char* ver, size_t nLen);

// Send Smartset command
ELOMT_API int EloSendSmartsetCommand( PSMARTSETDATA pSmartsetData );

// Get Smartset response
ELOMT_API int EloGetSmartsetResponse( PSMARTSETDATA pSmartsetData );

// Get Smartset response status 
ELOMT_API int EloGetSmartsetResponseStatus( PSMARTSETDATA pSmartsetData );

// Apply transactional Smartset command to retrieve the response data
//ELOMT_API bool EloSmartsetTransaction(SMARTSETDATA* pSS, SMARTSET_PKT* pSSPACKET=NULL, int nNumPkts=0, int* pnRespPkts=NULL, int nDelay=20);

// Clear pending GetPoint ioctl
ELOMT_API bool EloClearGetPoint (int nScreenIndex);

// Get multiple touch points from a touch screen
ELOMT_API bool EloGetMultiTouch (int nScreenIndex, MT_TOUCH* pMtTouch);

// Get Max Touch point of given screen index
ELOMT_API int EloGetMaxTouch (int nScreenIndex);

// Set Max Touch point of given screen index
ELOMT_API bool EloSetMaxTouch (int nScreenIndex, int nMaxTouch);

// Get Mouse Mode of given screen index
ELOMT_API bool EloIsForceMouse (int nScreenIndex);

// Get mouse mode
ELOMT_API int EloGetMouseMode (int nScreenIndex);

// Set mouse mode
ELOMT_API bool EloSetMouseMode (int nScreenIndex, int nMouseMode);

// Locking mechanism to protect share object in Smartset command
// EloAcquireSmartsetLock and EloReleaseSmartsetLock are used if you want to control
// Smartset transaction EloSendSmartsetCommand and EloGetSmartsetResponse. We strongly
// recommend to use EloSmartsetTransaction since it is thread safe. Acquire Smartset 
// lock prevent from other Smartset command to enter controller
//ELOMT_API bool EloAcquireSmartsetLock( int nScreenIndex, int nRetry = 1 );

// Release Smartset lock
ELOMT_API bool EloReleaseSmartsetLock( int nScreenIndex );

ELOMT_API bool EloSetBeepOptions (int nScreenIndex, const ELO_BEEP* pBeepOptions);
ELOMT_API bool EloGetBeepOptions (int nScreenIndex, ELO_BEEP* pBeepOptions);

// Get MouseExtraInformation of given screen index
ELOMT_API int EloGetMouseExtraInfo (int nScreenIndex);

// Set MouseExtraInformation of given screen index
ELOMT_API bool EloSetMouseExtraInfo (int nScreenIndex, int extraInfo);

ELOMT_API bool EloGetControllerSN (int nScreenIndex, wchar_t* pwszBuffer, int nBufLen);

ELOMT_API ULONGLONG EloGetPhysicalTouchCount (int nScreenIndex);
ELOMT_API bool EloSetPhysicalTouchCount (int nScreenIndex, ULONGLONG ullCount);

ELOMT_API bool EloIsEmulationMode (int nScreenIndex);

#ifdef __cplusplus
}
#endif


