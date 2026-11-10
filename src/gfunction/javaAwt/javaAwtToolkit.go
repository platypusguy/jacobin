/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaAwt

import "jacobin/src/gfunction/ghelpers"

func Load_Awt_Toolkit() {
	methods := map[string]int{
		"java/awt/Toolkit.<init>()V":                                                                                          0,
		"java/awt/Toolkit.loadSystemColors([I)V":                                                                              1,
		"java/awt/Toolkit.setDynamicLayout(Z)V":                                                                               1,
		"java/awt/Toolkit.isDynamicLayoutSet()Z":                                                                              0,
		"java/awt/Toolkit.isDynamicLayoutActive()Z":                                                                           0,
		"java/awt/Toolkit.getScreenSize()Ljava/awt/Dimension;":                                                                0,
		"java/awt/Toolkit.getScreenResolution()I":                                                                             0,
		"java/awt/Toolkit.getScreenInsets(Ljava/awt/GraphicsConfiguration;)Ljava/awt/Insets;":                                 1,
		"java/awt/Toolkit.getColorModel()Ljava/awt/image/ColorModel;":                                                         0,
		"java/awt/Toolkit.getFontList()[Ljava/lang/String;":                                                                   0,
		"java/awt/Toolkit.getFontMetrics(Ljava/awt/Font;)Ljava/awt/FontMetrics;":                                               1,
		"java/awt/Toolkit.sync()V":                                                                                            0,
		"java/awt/Toolkit.initAssistiveTechnologies()V":                                                                        0,
		"java/awt/Toolkit.newAWTError(Ljava/lang/Throwable;Ljava/lang/String;)V":                                               2,
		"java/awt/Toolkit.fallbackToLoadClassForAT(Ljava/lang/String;)V":                                                       1,
		"java/awt/Toolkit.loadAssistiveTechnologies()V":                                                                        0,
		"java/awt/Toolkit.getDefaultToolkit()Ljava/awt/Toolkit;":                                                              0,
		"java/awt/Toolkit.getImage(Ljava/lang/String;)Ljava/awt/Image;":                                                       1,
		"java/awt/Toolkit.getImage(Ljava/net/URL;)Ljava/awt/Image;":                                                          1,
		"java/awt/Toolkit.createImage(Ljava/lang/String;)Ljava/awt/Image;":                                                    1,
		"java/awt/Toolkit.createImage(Ljava/net/URL;)Ljava/awt/Image;":                                                       1,
		"java/awt/Toolkit.prepareImage(Ljava/awt/Image;IILjava/awt/image/ImageObserver;)Z":                                    4,
		"java/awt/Toolkit.checkImage(Ljava/awt/Image;IILjava/awt/image/ImageObserver;)I":                                       4,
		"java/awt/Toolkit.createImage(Ljava/awt/image/ImageProducer;)Ljava/awt/Image;":                                         1,
		"java/awt/Toolkit.createImage([B)Ljava/awt/Image;":                                                                     1,
		"java/awt/Toolkit.createImage([BII)Ljava/awt/Image;":                                                                  3,
		"java/awt/Toolkit.getPrintJob(Ljava/awt/Frame;Ljava/lang/String;Ljava/util/Properties;)Ljava/awt/PrintJob;":           3,
		"java/awt/Toolkit.getPrintJob(Ljava/awt/Frame;Ljava/lang/String;Ljava/awt/JobAttributes;Ljava/awt/PageAttributes;)Ljava/awt/PrintJob;": 4,
		"java/awt/Toolkit.beep()V":                                                                                             0,
		"java/awt/Toolkit.getSystemClipboard()Ljava/awt/datatransfer/Clipboard;":                                               0,
		"java/awt/Toolkit.getSystemSelection()Ljava/awt/datatransfer/Clipboard;":                                               0,
		"java/awt/Toolkit.getMenuShortcutKeyMask()I":                                                                           0,
		"java/awt/Toolkit.getMenuShortcutKeyMaskEx()I":                                                                         0,
		"java/awt/Toolkit.getLockingKeyState(I)Z":                                                                              1,
		"java/awt/Toolkit.setLockingKeyState(IZ)V":                                                                            2,
		"java/awt/Toolkit.getNativeContainer(Ljava/awt/Component;)Ljava/awt/Container;":                                        1,
		"java/awt/Toolkit.createCustomCursor(Ljava/awt/Image;Ljava/awt/Point;Ljava/lang/String;)Ljava/awt/Cursor;":             3,
		"java/awt/Toolkit.getBestCursorSize(II)Ljava/awt/Dimension;":                                                          2,
		"java/awt/Toolkit.getMaximumCursorColors()I":                                                                           0,
		"java/awt/Toolkit.isFrameStateSupported(I)Z":                                                                           1,
		"java/awt/Toolkit.setPlatformResources(Ljava/util/ResourceBundle;)V":                                                   1,
		"java/awt/Toolkit.initIDs()V":                                                                                          0,
		"java/awt/Toolkit.loadLibraries()V":                                                                                    0,
		"java/awt/Toolkit.getProperty(Ljava/lang/String;Ljava/lang/String;)Ljava/lang/String;":                                 2,
		"java/awt/Toolkit.getSystemEventQueue()Ljava/awt/EventQueue;":                                                         0,
		"java/awt/Toolkit.getSystemEventQueueImpl()Ljava/awt/EventQueue;":                                                     0,
		"java/awt/Toolkit.getEventQueue()Ljava/awt/EventQueue;":                                                               0,
		"java/awt/Toolkit.createDragGestureRecognizer(Ljava/lang/Class;Ljava/awt/dnd/DragSource;Ljava/awt/Component;ILjava/awt/dnd/DragGestureListener;)Ljava/awt/dnd/DragGestureRecognizer;": 5,
		"java/awt/Toolkit.getDesktopProperty(Ljava/lang/String;)Ljava/lang/Object;":                                            1,
		"java/awt/Toolkit.setDesktopProperty(Ljava/lang/String;Ljava/lang/Object;)V":                                           2,
		"java/awt/Toolkit.lazilyLoadDesktopProperty(Ljava/lang/String;)Ljava/lang/Object;":                                      1,
		"java/awt/Toolkit.initializeDesktopProperties()V":                                                                      0,
		"java/awt/Toolkit.addPropertyChangeListener(Ljava/lang/String;Ljava/beans/PropertyChangeListener;)V":                   2,
		"java/awt/Toolkit.removePropertyChangeListener(Ljava/lang/String;Ljava/beans/PropertyChangeListener;)V":                2,
		"java/awt/Toolkit.getPropertyChangeListeners()[Ljava/beans/PropertyChangeListener;":                                    0,
		"java/awt/Toolkit.getPropertyChangeListeners(Ljava/lang/String;)[Ljava/beans/PropertyChangeListener;":                 1,
		"java/awt/Toolkit.isAlwaysOnTopSupported()Z":                                                                           0,
		"java/awt/Toolkit.isModalityTypeSupported(Ljava/awt/Dialog$ModalityType;)Z":                                            1,
		"java/awt/Toolkit.isModalExclusionTypeSupported(Ljava/awt/Dialog$ModalExclusionType;)Z":                                1,
		"java/awt/Toolkit.deProxyAWTEventListener(Ljava/awt/event/AWTEventListener;)Ljava/awt/event/AWTEventListener;":        1,
		"java/awt/Toolkit.addAWTEventListener(Ljava/awt/event/AWTEventListener;J)V":                                            2,
		"java/awt/Toolkit.removeAWTEventListener(Ljava/awt/event/AWTEventListener;)V":                                          1,
		"java/awt/Toolkit.enabledOnToolkit(J)Z":                                                                                1,
		"java/awt/Toolkit.countAWTEventListeners(J)I":                                                                          1,
		"java/awt/Toolkit.getAWTEventListeners()[Ljava/awt/event/AWTEventListener;":                                            0,
		"java/awt/Toolkit.getAWTEventListeners(J)[Ljava/awt/event/AWTEventListener;":                                          1,
		"java/awt/Toolkit.notifyAWTEventListeners(Ljava/awt/AWTEvent;)V":                                                       1,
		"java/awt/Toolkit.mapInputMethodHighlight(Ljava/awt/im/InputMethodHighlight;)Ljava/util/Map;":                          1,
		"java/awt/Toolkit.createPropertyChangeSupport(Ljava/awt/Toolkit;)Ljava/beans/PropertyChangeSupport;":                 1,
		"java/awt/Toolkit.areExtraMouseButtonsEnabled()Z":                                                                      0,
		"java/awt/Toolkit.lambda$loadAssistiveTechnologies$1(Ljava/util/Map;Ljava/lang/String;)Z":                              2,
		"java/awt/Toolkit.lambda$loadAssistiveTechnologies$0(Ljava/lang/ClassLoader;Ljava/util/Set;Ljava/util/Map;)Ljava/lang/Void;": 3,
	}

	for signature, paramSlots := range methods {
		ghelpers.MethodSignatures[signature] = ghelpers.GMeth{
			ParamSlots: paramSlots,
			GFunction:  ghelpers.TrapFunction,
		}
	}

	ghelpers.MethodSignatures["java/awt/Toolkit.<clinit>()V"] = ghelpers.GMeth{
		ParamSlots: 0,
		GFunction:  ghelpers.ClinitGeneric,
	}

	ghelpers.MethodSignatures["java/awt/Toolkit.initStatic()V"] = ghelpers.GMeth{
		ParamSlots: 0,
		GFunction:  ghelpers.JustReturn,
	}
}