# Android Build Blockers

## Current blocker

Running Gradle tasks in this environment currently fails with a native loader error for `libnative-platform.so`.

JNI/CMake wiring for `android/corebridge` is now present (`externalNativeBuild` + `CMakeLists.txt`), but task execution remains blocked by the host environment issue above.

## Wrapper status in this scaffold

This stage of the scaffold does not include the full Gradle wrapper artifacts (`gradle/wrapper/gradle-wrapper.jar` and `gradle/wrapper/gradle-wrapper.properties`).

To keep task execution unblocked for local environments that already have Gradle installed, `./gradlew` is currently a lightweight shim that delegates directly to the installed `gradle` binary.
