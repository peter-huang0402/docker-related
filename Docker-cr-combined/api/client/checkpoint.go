// +build experimental
// ./docker/api/client/checkpoint.go     //peter
package client

import (
	"fmt"
	"github.com/docker/docker/api/types"
	Cli "github.com/docker/docker/cli"
	flag "github.com/docker/docker/pkg/mflag"
)

// CmdCheckpoint checkpoints the process running in a container
//
// Usage: docker checkpoint CONTAINER
func (cli *DockerCli) CmdCheckpoint(args ...string) error {
	cmd := Cli.Subcmd("checkpoint", []string{"CONTAINER"}, Cli.DockerCommands["checkpoint"].Description, true)
	cmd.Require(flag.Min, 1)

	var (
		flImgDir       = cmd.String([]string{"-image-dir"}, "", "directory for storing checkpoint image files")
		flWorkDir      = cmd.String([]string{"-work-dir"}, "", "directory for storing log file")
		flLeaveRunning = cmd.Bool([]string{"-leave-running"}, false, "leave the container running after checkpoint")
                flPrevImgDir      = cmd.String([]string{"-prev-image-dir"}, "", "previous Image directory")  //peter
		flTrackMem = cmd.Bool([]string{"-track-mem"}, false, "Enable track memory flag")   // peter
                flEnablePreDump = cmd.Bool([]string{"-predump"}, false, "Enable predump flag")   // peter 
                 flLazyPages = cmd.Bool([]string{"-lazy-pages"}, false, "Enable lazy-pages flag for checkpoint")   // peter for lazy-pages
                 flLazyPagesPort = cmd.Int([]string{"-lazy-pages-port"}, 0 , "Port of Lazy-Pages used with lazy-pages flag")   // peter for lazy-pages 
                flAutoDedup = cmd.Bool([]string{"-auto-dedup"}, false, "Enable auto-dedup flag")   // peter 
                flPageServer = cmd.Bool([]string{"-page-server"}, false, "Enable Page-Server flag for RPC")   // peter 
                flAddress       = cmd.String([]string{"-address"}, "", "IP address of Page-Server used with Page-Server flag")    //peter 
                flPort      = cmd.Int([]string{"-port"},0 , "Port of Page-Server used with Page-Server flag")  
	)

	if err := cmd.ParseFlags(args, true); err != nil {
		return err
	}

	if cmd.NArg() < 1 {
		cmd.Usage()
		return nil
	}

      

        // add by peter
        var p int32
        var tmp int = *flPort;
         fmt.Printf("flPageServer:[%v], flAddress:[%s], flPort:[%v] \n", *flPageServer, *flAddress, tmp )  // peter

        if  *flPageServer == true {
              if *flAddress =="" {
                   return fmt.Errorf("Page Server enabled but address is not assigned");
              }
              if *flPort ==0 {
                  return fmt.Errorf("Page Server enabled but port is not assigned");
              } 
              
              p = int32(tmp)
        }

        // add by peter lazy-pages 
        var p2 int32 
        var tmp2 int = *flLazyPagesPort;
        if  *flLazyPages == true {
               if *flLazyPagesPort ==0 {
                  return fmt.Errorf("LazyPage Server enabled but lazy-pages-port is not assigned");
              }            
               p2 = int32(tmp2)
         } 




	criuOpts := types.CriuConfig{
		ImagesDirectory: *flImgDir,
		WorkDirectory:   *flWorkDir,
		LeaveRunning:    *flLeaveRunning,
                PrevImagesDirectory:  *flPrevImgDir,  //peter 
                TrackMemory:  *flTrackMem,  //peter
                LazyPages:     *flLazyPages,    //peter for lazy-pages
                LazyPagesPort:     p2,          //peter for lazy-pages
                EnablePreDump:  *flEnablePreDump, //peter
                AutoDedup:      *flAutoDedup,     //peter 
                PageServer:     *flPageServer,    //peter 
                Address:        *flAddress,       //peter 
                Port:           p,         //peter 
	}

	var encounteredError error
	for _, name := range cmd.Args() {
		err := cli.client.ContainerCheckpoint(name, criuOpts)
		if err != nil {
			fmt.Fprintf(cli.err, "%s\n", err)
			encounteredError = fmt.Errorf("Error: failed to checkpoint one or more containers")
		} else {
			fmt.Fprintf(cli.out, "%s\n", name)
		}
	}
	return encounteredError
}
