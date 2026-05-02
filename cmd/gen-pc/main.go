// gen-pc detects the active ROS 2 installation and writes a pkg-config
// file (rcl-go.pc) that exposes all the cgo flags needed to build this module.
//
// Usage:
//
//	go run github.com/lesomnus/rcl/cmd/gen-pc -o /path/to/rcl-go.pc
//	# then add that directory to PKG_CONFIG_PATH
//
// If -o is omitted the file is written to the current directory.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// known_distros lists ROS 2 distributions in preference order (newest first).
var known_distros = []string{
	"rolling",
	"jazzy",
	"iron",
	"humble",
	"galactic",
	"foxy",
}

// include_packages lists ROS 2 sub-packages whose headers are needed.
var include_packages = []string{
	"action_msgs",
	"builtin_interfaces",
	"rcl",
	"rcutils",
	"rcl_action",
	"rcl_yaml_param_parser",
	"rmw",
	"rosidl_dynamic_typesupport",
	"rosidl_runtime_c",
	"rosidl_typesupport_interface",
	"rosidl_typesupport_introspection_c",
	"service_msgs",
	"type_description_interfaces",
	"unique_identifier_msgs",
}

var libs = []string{
	"rcl",
	"rmw",
	"rosidl_runtime_c",
	"rosidl_typesupport_c",
	"rcutils",
	"rcl_action",
	"rmw_implementation",
}

func hasRCL(prefix string) bool {
	_, err := os.Stat(filepath.Join(prefix, "include", "rcl"))
	return err == nil
}

func detectRosPrefix() (string, error) {
	// 1. AMENT_PREFIX_PATH (set by `source /opt/ros/<distro>/setup.bash`).
	//    May be colon-separated; pick the first entry that contains rcl.
	if raw := os.Getenv("AMENT_PREFIX_PATH"); raw != "" {
		for _, prefix := range strings.Split(raw, ":") {
			if hasRCL(prefix) {
				return prefix, nil
			}
		}
	}

	// 2. ROS_DISTRO (also set by the setup script).
	if distro := os.Getenv("ROS_DISTRO"); distro != "" {
		prefix := filepath.Join("/opt/ros", distro)
		if hasRCL(prefix) {
			return prefix, nil
		}
	}

	// 3. Scan /opt/ros/, preferring known distros (newest first).
	entries, err := os.ReadDir("/opt/ros")
	if err == nil {
		installed := map[string]bool{}
		for _, e := range entries {
			if e.IsDir() {
				installed[e.Name()] = true
			}
		}
		for _, distro := range known_distros {
			if installed[distro] {
				prefix := filepath.Join("/opt/ros", distro)
				if hasRCL(prefix) {
					return prefix, nil
				}
			}
		}
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			prefix := filepath.Join("/opt/ros", e.Name())
			if hasRCL(prefix) {
				return prefix, nil
			}
		}
	}

	return "", fmt.Errorf("ROS 2 installation not found; " +
		"source /opt/ros/<distro>/setup.bash or set AMENT_PREFIX_PATH / ROS_DISTRO")
}

func main() {
	out := flag.String("o", "-", "output path for the generated .pc file")
	flag.Parse()

	prefix, err := detectRosPrefix()
	if err != nil {
		fmt.Fprintln(os.Stderr, "gen-pc:", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "gen-pc: using ROS prefix %s\n", prefix)

	var cflags []string
	for _, pkg := range include_packages {
		cflags = append(cflags, fmt.Sprintf("-I%s/include/%s", prefix, pkg))
	}

	var ldflags []string
	ldflags = append(ldflags,
		"-L/usr/lib", "-Wl,-rpath=/usr/lib",
		fmt.Sprintf("-L%s/lib", prefix),
		fmt.Sprintf("-Wl,-rpath=%s/lib", prefix),
	)
	for _, lib := range libs {
		ldflags = append(ldflags, "-l"+lib)
	}

	content := fmt.Sprintf(`Name: rcl-go
Description: cgo flags for github.com/lesomnus/rcl (ROS 2 prefix: %s)
Version: 0.0.0

Cflags: %s
Libs: %s
`, prefix, strings.Join(cflags, " "), strings.Join(ldflags, " "))

	if *out == "-" {
		fmt.Fprint(os.Stdout, content)
	} else {
		if err := os.MkdirAll(filepath.Dir(*out), 0755); err != nil {
			fmt.Fprintln(os.Stderr, "gen-pc:", err)
			os.Exit(1)
		}
		if err := os.WriteFile(*out, []byte(content), 0644); err != nil {
			fmt.Fprintln(os.Stderr, "gen-pc:", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "gen-pc: wrote %s\n", *out)
	}
}
